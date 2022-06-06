package ws

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gocs/map/models"
	"github.com/gocs/map/store"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == os.Getenv("GOCS_MAP_ORIGIN")
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	filename := "./static/land.json"

	node := &models.NodeRecv{}
	b := &bytes.Buffer{}
	for {
		if err := c.conn.ReadJSON(node); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
				break
			}
			c.conn.WriteMessage(websocket.CloseAbnormalClosure, []byte(err.Error()))
			break
		}

		var fn func(m store.Map) (store.Map, error)
		switch node.Action {
		case "add":
			fn = models.AddMap(*node)
		case "del":
			fn = models.DelMap(*node)
		case "set":
			fn = models.SetMap(*node)
		default:
			log.Println("err missing fn:")
			continue
		}

		if err := store.MapMap(filename, fn); err != nil {
			log.Println("err MapMap:", err)
			continue
		}

		if err := json.NewEncoder(b).Encode(&map[string]string{"action": node.Action}); err != nil {
			log.Println("err Encode:", err)
			continue
		}
		c.hub.broadcast <- b.Bytes()
		b.Reset()
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.conn.WriteMessage(websocket.CloseAbnormalClosure, []byte(err.Error()))
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				c.conn.WriteMessage(websocket.CloseAbnormalClosure, []byte(err.Error()))
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.conn.WriteMessage(websocket.CloseAbnormalClosure, []byte(err.Error()))
				return
			}
		}
	}
}
