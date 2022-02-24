package models

import (
	"github.com/gocs/map/store"
)

type NodeRecv struct {
	Action string `json:"action"`
	ID     int    `json:"id"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Type   string `json:"type"`
}

// AddMap adds new node
// does not not realign missing ids
func AddMap(node NodeRecv) func(m store.Map) (store.Map, error) {
	return func(m store.Map) (store.Map, error) {
		highestID := 0
		if len(m.Nodes) > 0 {
			for _, n := range m.Nodes {
				if n.ID > highestID {
					highestID = n.ID
				}
			}
		}

		m.Nodes = append(m.Nodes, store.Node{
			ID:   highestID + 1,
			X:    node.X,
			Y:    node.Y,
			Type: node.Type,
		})
		return m, nil
	}
}

// DelMap only uses node.ID and WILL ignore the other fields
// retains the id numbers; there could have missing numbers
func DelMap(node NodeRecv) func(m store.Map) (store.Map, error) {
	return func(m store.Map) (store.Map, error) {
		m.Nodes = filter(m.Nodes, func(n store.Node) bool {
			return n.ID != node.ID
		})
		return m, nil
	}
}

func filter(ss []store.Node, h func(store.Node) bool) []store.Node {
	ret := []store.Node{}
	for _, s := range ss {
		if h(s) {
			ret = append(ret, s)
		}
	}
	return ret
}

func SetMap(node NodeRecv) func(m store.Map) (store.Map, error) {
	return func(m store.Map) (store.Map, error) {
		m.Nodes[node.ID] = store.Node{
			ID:   node.ID,
			X:    node.X,
			Y:    node.Y,
			Type: node.Type,
		}
		return m, nil
	}
}
