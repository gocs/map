package main

import (
	"log"
	"net/http"

	"github.com/gocs/map/ws"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/updaterws", ws.UpdaterWS())
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/xqcLewd.ico")
	})

	log.Println("Listening on :3131...")
	if err := http.ListenAndServe(":3131", nil); err != nil {
		log.Fatal(err)
	}
}
