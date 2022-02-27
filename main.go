package main

import (
	"log"
	"net/http"

	"github.com/gocs/map/ws"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/", fs)
	http.HandleFunc("/updaterws", ws.UpdaterWS())
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/xqcLewd.ico")
	})

	// my local 3000 is buzy 😂
	log.Println("Listening on :3131...")
	if err := http.ListenAndServe(":3131", nil); err != nil {
		log.Fatal(err)
	}
}
