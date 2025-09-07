package server

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed sylgraph/dist/*
var sylgraph_statics embed.FS

func GetStaticServer() (http.Handler, error) {
	fsys, err := fs.Sub(sylgraph_statics, "sylgraph/dist")
	if err != nil {
		slog.Error("failed to get static filesystem")
		return nil, err
	}

	return http.FileServer(http.FS(fsys)), nil
}

func (server *Server) HandleStatic(w http.ResponseWriter, r *http.Request) {
	if server == nil {
		return
	}
	// s := server

	path := filepath.Clean(r.URL.Path)
	if path == "/" { // Add other paths that you route on the UI side here
		path = "index.html"
	}
	path = strings.TrimPrefix(path, "/")

	slog.Info("urrrrrrrrrrrrrrrrr " + r.URL.Path)

	type Resp struct {
		Hi string `json:"hi"`
	}

	WriteJson(Resp{
		Hi: "catchall",
	}, w)
}
