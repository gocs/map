package ws

import (
	"log"
	"net/http"
)

func UpdaterWS() http.HandlerFunc {
	hub := newHub()
	go hub.run()

	return func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
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
