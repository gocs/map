package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Node struct {
	ID   int    `json:"id"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Type string `json:"type"`
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

func MapMap(filename string, handler func(m Map) (Map, error)) error {
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
