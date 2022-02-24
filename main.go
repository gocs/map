package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocs/map/store"
	"github.com/gocs/map/ws"
)

func main() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/landjson", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		var n store.Node

		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			log.Println("err Decode:", err)
			return
		}
		err := store.MapMap("./static/land.json", func(m store.Map) (store.Map, error) {
			m.Nodes[n.ID] = n
			return m, nil
		})
		if err != nil {
			log.Println("err mapMap:", err)
			return
		}
	})

	http.HandleFunc("/appendcell", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		var n store.Node
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			log.Println("err Decode:", err)
			return
		}
		err := store.MapMap("./static/land.json", func(m store.Map) (store.Map, error) {
			n.ID = len(m.Nodes)
			m.Nodes = append(m.Nodes, n)
			return m, nil
		})
		if err != nil {
			log.Println("err mapMap:", err)
			return
		}
	})

	http.HandleFunc("/deletecell", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			return
		}

		var d store.Deleter
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			log.Println("err Decode:", err)
			return
		}
		err := store.MapMap("./static/land.json", func(m store.Map) (store.Map, error) {
			m.Nodes = append(m.Nodes[:d.NodeID], m.Nodes[d.NodeID+1:]...)
			for i := 0; i < len(m.Nodes); i++ {
				m.Nodes[i].ID = i
			}
			return m, nil
		})
		if err != nil {
			log.Println("err mapMap:", err)
			return
		}
	})

	http.HandleFunc("/updaterws", ws.UpdaterWS())

	// my local 3000 is buzy 😂
	log.Println("Listening on :3131...")
	if err := http.ListenAndServe(":3131", nil); err != nil {
		log.Fatal(err)
	}
}
