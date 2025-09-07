package server

import (
	"log/slog"
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

	nodeIdRefsMap := map[int][]lsp.Location{}
	for gTarget, v := range s.store.GLinkStore {
		if len(v.Defs) > 0 {
			for _, def := range v.Defs {
				target, _ := data.GetFileGTarget(def.URI)
				node := server.graphStore.StoreAndGetId(target, 0, NodeKindFile)
				refs, found := nodeIdRefsMap[node.Id]
				if found {
					nodeIdRefsMap[node.Id] = append(refs, v.Refs...)
				} else {
					refs = []lsp.Location{}
					nodeIdRefsMap[node.Id] = refs
				}
			}
		} else {
			_target, _, _ := gTarget.SplitHeading()
			node := server.graphStore.StoreAndGetId(string(_target), 0, NodeKindUnresolvedFile)
			refs, found := nodeIdRefsMap[node.Id]
			if found {
				nodeIdRefsMap[node.Id] = append(refs, v.Refs...)
			} else {
				refs = []lsp.Location{}
				nodeIdRefsMap[node.Id] = refs
			}
		}
	}
	// add tags
	for tag, refs := range s.store.Tags {
		node := server.graphStore.StoreAndGetId(string(tag), 0, NodeKindTag)
		nodeIdRefsMap[node.Id] = refs
	}

	g := NewGraph()
	maxCon := 0
	minCon := 99999

	linkMap := map[int]map[int]bool{}
	for nodeId, targets := range nodeIdRefsMap {

		connections := len(targets)

		minCon = min(minCon, connections)
		maxCon = max(maxCon, connections)

		node, _ := s.graphStore.GetNodeFromId(nodeId)
		g.Nodes = append(g.Nodes, node)

		for _, target := range targets {
			target, _ := data.GetFileGTarget(target.URI)

			sourceNode, found := s.graphStore.GetNodeFromName(string(target))
			if !found {
				slog.Info("Node not found should have been there")
			}
			sourceId := sourceNode.Id
			_, found = linkMap[sourceId]
			if found {
				linkMap[sourceId][nodeId] = true
			} else {
				linkMap[sourceId] = map[int]bool{}
				linkMap[sourceId][nodeId] = true
			}

		}

	}

	for sourceId, targetMap := range linkMap {
		for targetId := range targetMap {
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
