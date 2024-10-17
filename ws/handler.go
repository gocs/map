package ws

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

func UpdaterWS(expected string) http.HandlerFunc {
	hub := newHub()
	go hub.run()

	var upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			slog.Info("check origin", "actual", origin, "expected", expected)
			return origin == expected
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		serveWs(upgrader, hub, w, r)
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(u *websocket.Upgrader, hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
