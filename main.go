package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gocs/map/static"
	"github.com/gocs/map/ws"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	mapOrigin := os.Getenv("GOCS_MAP_ORIGIN")
	if mapOrigin == "" {
		mapOrigin = "http://localhost:" + port
	}

	a, err := static.Assets()
	if err != nil {
		log.Fatalln("assets err:", err)
		return
	}

	http.Handle("/", a)
	http.HandleFunc("/land.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "land.json")
	})
	http.HandleFunc("/updaterws", ws.UpdaterWS(mapOrigin))
	http.HandleFunc("/ws", ws.UpdaterWS(mapOrigin))

	log.Printf("Listening on :%s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
