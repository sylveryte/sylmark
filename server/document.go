package server

import (
	"log/slog"
	"net/http"
	"sylmark-server/data"
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

	node, found := server.graphStore.GetNodeFromId(py.Id)
	if found {
		locs := server.store.GetGTargetDefinition(data.GTarget(node.Name))
		if len(locs) > 0 {
			loc := locs[0]
			server.showDocument(loc.URI, loc.Range)
		}
	}
}
