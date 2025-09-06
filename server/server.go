package server

import (
	"log/slog"
	"net/http"
	"sylmark-server/data"
	"sylmark-server/lsp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	graphStore   *GraphStore
	store        *data.Store
	rootPath     string
	Config       *data.Config
	showDocument lsp.ShowDocumentFx
}

func NewServer(store *data.Store, config *data.Config, showDocument lsp.ShowDocumentFx) (server *Server) {
	return &Server{
		store:        store,
		graphStore:   newGraphStore(),
		Config:       config,
		showDocument: showDocument,
	}
}

func (s *Server) StartAndListen() {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*", "http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Logger)
	v1 := chi.NewRouter()
	r.Mount("/v1", v1)
	s.SetupRoutes(v1)
	slog.Info("Staring server at 7462")
	//port 293001-293010
	http.ListenAndServe(":7462", r)
}
