package server

import (
	"net/http"
	"sylmark/data"
	"sylmark/lsp"
)

type Node struct {
	Id   int      `json:"id"`
	Name string   `json:"name"`
	Val  int      `json:"val"`
	Kind NodeKind `json:"kind"`
	uri  lsp.DocumentURI
}
type NodeKind int

const (
	NodeKindFile           NodeKind = 1
	NodeKindTag            NodeKind = 2
	NodeKindUnresolvedFile NodeKind = 3
)

type Link struct {
	Source int `json:"source"`
	Target int `json:"target"`
}
type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

func NewGraph() Graph {
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

	graphMap := map[string][]lsp.Location{}
	for _, v := range s.store.GLinkStore {
		for _, def := range v.Defs {
			node, _ := data.GetFileGTarget(def.URI)
			nodeId := string(node)
			refs, found := graphMap[nodeId]
			if !found {
				refs = []lsp.Location{}
				graphMap[nodeId] = refs
			} else {
				graphMap[nodeId] = append(refs, v.Refs...)
			}
		}
	}

	g := NewGraph()
	maxCon := 0
	minCon := 99999

	linkMap := map[int]map[int]bool{}
	for nodeId, targets := range graphMap {
		connections := len(targets)

		minCon = min(minCon, connections)
		maxCon = max(maxCon, connections)

		node := s.graphStore.StoreAndGetId(nodeId, 0, NodeKindFile)
		g.Nodes = append(g.Nodes, node)

		for _, target := range targets {
			target, _ := data.GetFileGTarget(target.URI)

			targetNode, _ := s.graphStore.GetNodeFromName(nodeId)
			sourceNode, _ := s.graphStore.GetNodeFromName(string(target))
			_, found := linkMap[sourceNode.Id]
			if found {
				linkMap[sourceNode.Id][targetNode.Id] = true
			} else {
				linkMap[sourceNode.Id] = map[int]bool{}
				linkMap[sourceNode.Id][targetNode.Id] = true
			}

		}
	}

	for sourceId, targetMap := range linkMap {
		for targetId, _ := range targetMap {
			g.Links = append(g.Links, Link{
				Source: sourceId,
				Target: targetId,
			})

		}

	}

	// updated nodes with better size
	nodes := []Node{}
	for _, n := range g.Nodes {
		// slog.Info(fmt.Sprintf("%d => (%d > %d) == %d", n.Val, maxCon, minCon, getSize(n.Val, maxCon, minCon)))
		connections, found := linkMap[n.Id]
		var totalConnections int
		if !found {
			totalConnections = 0
		} else {
			totalConnections = len(connections)
		}
		n.Val = getSize(totalConnections, maxCon, minCon)
		nodes = append(nodes, n)
	}
	g.Nodes = nodes

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
