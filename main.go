package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gocs/map/ws"
)

type Node struct {
	NodeID int     `json:"id"`
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Type   string  `json:"type"`
}

type Circle struct {
	Polygon     int `json:"polygon"`
	Radius      int `json:"radius"`
	StrokeWidth int `json:"stroke_width"`
}

type Colors struct {
	Cliff  string `json:"cliff"`
	Forest string `json:"forest"`
	Shore  string `json:"shore"`
	Sea    string `json:"sea"`
}

type Options struct {
	Colors Colors `json:"colors"`
}

type Line struct {
	Color       string `json:"color"`
	StrokeWidth int    `json:"stroke_width"`
}

type Deleter struct {
	NodeID int `json:"id"`
}

type Map struct {
	Nodes   []Node  `json:"nodes"`
	Circle  Circle  `json:"circle"`
	Options Options `json:"options"`
	Line    Line    `json:"line"`
}

func main() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/landjson", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		var n Node

		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			log.Println("err Decode:", err)
			return
		}
		manageMap("./static/land.json", func(m Map) (Map, error) {
			m.Nodes[n.NodeID] = n
			return m, nil
		})
	})

	http.HandleFunc("/appendcell", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		var n Node
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			log.Println("err Decode:", err)
			return
		}
		manageMap("./static/land.json", func(m Map) (Map, error) {
			n.NodeID = len(m.Nodes)
			m.Nodes = append(m.Nodes, n)
			return m, nil
		})
	})

	http.HandleFunc("/deletecell", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			return
		}

		var d Deleter
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			log.Println("err Decode:", err)
			return
		}
		manageMap("./static/land.json", func(m Map) (Map, error) {
			m.Nodes = append(m.Nodes[:d.NodeID], m.Nodes[d.NodeID+1:]...)
			for i := 0; i < len(m.Nodes); i++ {
				m.Nodes[i].NodeID = i
			}
			return m, nil
		})
	})

	http.HandleFunc("/updaterws", ws.UpdaterWS())

	// my local 3000 is buzy 😂
	log.Println("Listening on :3131...")
	if err := http.ListenAndServe(":3131", nil); err != nil {
		log.Fatal(err)
	}
}

func manageMap(filename string, handler func(m Map) (Map, error)) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	var m Map
	if err := json.Unmarshal([]byte(byteValue), &m); err != nil {
		return err
	}
	m, err = handler(m)
	if err != nil {
		return err
	}
	ba, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, ba, 0644)
}
