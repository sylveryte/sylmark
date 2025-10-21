package server

import (
	"log/slog"
	"net/http"
	"sylmark/lsp"
)

type ShowDocumentParams struct {
	Id int `json:"id"`
}

func (server *Server) ShowDocument(w http.ResponseWriter, r *http.Request) {
	var py ShowDocumentParams
	err := ReadJSON(&py, r)
	if err != nil {
		slog.Error("Failed to read request body")
		return
	}

	node, found := server.graphStore.nodeStore.get(NodeId(py.Id))
	if found {
		uri, ok := server.store.GetUri(node.InternalId)
		if ok {
			server.showDocument(uri, false, lsp.Range{})
		}
	}
}
