package server

import (
	"log/slog"
	"strings"
	"sylmark/data"
)

// will be data.Id and fromlast +100 for tags
type NodeId uint
type Node struct {
	Id         NodeId   `json:"id"`
	InternalId data.Id  `json:"linkId"`
	Name       string   `json:"name"`
	Val        int      `json:"val"` // determines size
	Kind       NodeKind `json:"kind"`
	Path       string   `json:"path"` // relative path
}
type NodeKind int16

const (
	NodeKindFile           NodeKind = 1
	NodeKindTag            NodeKind = 2
	NodeKindUnresolvedFile NodeKind = 3
)

type NodeStore map[NodeId]Node

func newNodeStore() NodeStore {
	return NodeStore{}
}

func (store *NodeStore) add(node Node) NodeId {
	ns := *store
	id := node.Id
	if node.Id == 0 {
		id = NodeId(len(ns) + 100) // +100 so tags nodes can maintain safe distance
		node.Id = id
	}
	ns[id] = node
	return id
}
func (store *NodeStore) get(id NodeId) (Node, bool) {
	ns := *store
	n, ok := ns[id]
	return n, ok
}

func (store *NodeStore) updateVal(id NodeId, val int) {
	s := *store
	n, ok := s[id]
	if ok {
		n.Val = val
		s[id] = n
	}
}

type LinkStore map[NodeId]map[NodeId]int

func newLinkStore() LinkStore {
	return LinkStore{}
}

func (ls *LinkStore) add(source NodeId, target NodeId) {
	l := *ls
	tm, ok := l[source]
	if !ok {
		m := map[NodeId]int{}
		m[target] = 1
		l[source] = m
	} else {
		count, _ := tm[target]
		tm[target] = count + 1
		l[source] = tm
	}
}

type GraphStore struct {
	nodeStore NodeStore
	linkStore LinkStore
	minCon    int
	maxCon    int
}

func newGraphStore() *GraphStore {
	return &GraphStore{
		nodeStore: newNodeStore(),
		linkStore: newLinkStore(),
	}
}

func (server *Server) LoadGraph() {
	if server == nil && server.graphStore != nil {
		slog.Error("GraphStore is nil")
		return
	}
	s := server

	// store everything in NodeStore
	// adding resolved files
	for id, link := range s.store.LinkStore {

		// node
		uri, ok := s.store.GetUri(id)
		relPath, err := s.store.GetPathRelRoot(uri)
		if ok && len(uri) != 0 && err == nil {
			target, _ := data.GetTarget(uri)
			s.graphStore.nodeStore.add(Node{
				Id:         NodeId(id),
				InternalId: id,
				Name:       string(target),
				Kind:       NodeKindFile,
				Path:       relPath,
			})
		}

		// links
		for _, r := range link.Refs {
			for _, l := range r {
				s.graphStore.linkStore.add(NodeId(id), NodeId(l.Id))
			}
		}
	}
	// add shadows
	for id, tars := range s.store.IdStore.ShadowTargets {
		var t data.Target
		if len(tars) == 1 {
			t = tars[0]
		} else {
			for _, tar := range tars {
				if strings.ContainsRune(string(tar), '/') {
					continue
				}
				t = tar
			}
			if len(t) == 0 {
				t = tars[0]
			}
		}
		s.graphStore.nodeStore.add(Node{
			Id:         NodeId(id),
			InternalId: id,
			Name:       string(t),
			Kind:       NodeKindUnresolvedFile,
			Path:       string(t),
		})

	}

	// add tags
	for tag, refs := range s.store.Tags {
		nodeId := s.graphStore.nodeStore.add(
			Node{
				Name: string(tag),
				Kind: NodeKindTag,
			},
		)
		for _, l := range refs {
			id := s.store.GetIdFromURI(l.URI)
			s.graphStore.linkStore.add(nodeId, NodeId(id))
		}
	}

	s.graphStore.maxCon = 0
	s.graphStore.minCon = 99999

	// link counts update
	for id, tm := range s.graphStore.linkStore {
		con := 0
		for _, i := range tm {
			con += i
		}

		s.graphStore.minCon = min(s.graphStore.minCon, con)
		s.graphStore.maxCon = max(s.graphStore.maxCon, con)

		s.graphStore.nodeStore.updateVal(id, con)
		s.graphStore.nodeStore.updateVal(id, con)
	}
}
