package server

import (
	"net/http"
)

type Link struct {
	Source NodeId `json:"source"`
	Target NodeId `json:"target"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

func newGraph() Graph {
	return Graph{
		Nodes: []Node{},
		Links: []Link{},
	}
}

func (server *Server) GetGraph(w http.ResponseWriter, r *http.Request) {
	if server == nil {
		return
	}
	s := server

	// to refresh the data
	s.graphStore = newGraphStore()
	s.LoadGraph()

	g := newGraph()
	gs := s.graphStore

	for _, n := range gs.nodeStore {
		n.Val = getSize(n.Val, gs.maxCon, gs.minCon)
		g.Nodes = append(g.Nodes, n)
	}

	for sourceId, tm := range gs.linkStore {
		for targetId := range tm {
			g.Links = append(g.Links, Link{
				Source: sourceId,
				Target: targetId,
			})
		}
	}

	WriteJson(g, w)
}

// returns 1, 2, 3
func getSize(connections int, maxCon int, minCon int) int {

	if connections < 5 {
		return 3
	}

	normal := (connections - minCon) * 100 / (maxCon - minCon)

	if normal < 50 {
		return 4
	} else if normal < 75 {
		return 5
	}

	return 6
}
