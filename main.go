package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gocs/map/ws"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/updaterws", ws.UpdaterWS())
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/xqcLewd.ico")
	})

	log.Printf("Listening on :%s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
