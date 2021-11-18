package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

type Land struct {
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

		l, err := getLand()
		if err != nil {
			log.Println("err getLand:", err)
			return
		}

		var n Node

		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			log.Println("err Decode:", err)
			return
		}

		l.Nodes[n.NodeID] = n
		if err := json.NewEncoder(w).Encode(&n); err != nil {
			log.Println("err Encode:", err)
			return
		}

		if err := setLand(l); err != nil {
			log.Println("err setLand:", err)
			return
		}
	})

	log.Println("Listening on :3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

func getLand() (Land, error) {
	f, err := os.Open("./static/land.json")
	if err != nil {
		return Land{}, err
	}
	defer f.Close()

	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return Land{}, err
	}

	var l Land
	if err := json.Unmarshal([]byte(byteValue), &l); err != nil {
		return Land{}, err
	}
	return l, nil
}

func setLand(l Land) error {
	ba, err := json.Marshal(&l)
	if err != nil {
		return err
	}

	return os.WriteFile("./static/land.json", ba, 0644)
}
