package server

import (
	"log/slog"
	"net/http"
	"sylmark/data"
	"sylmark/lsp"

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

func (s *Server) StartAndListen() error {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:7462", "http://localhost:5173"},
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
	fsHandler, err := GetStaticServer()
	if err != nil {
		slog.Error("failed to get fileServerHandler " + err.Error())
		return err
	}
	r.Get("/*", fsHandler.ServeHTTP)
	port := "7462"
	slog.Info("Staring server at " + port)
	go http.ListenAndServe(":"+port, r)
	s.showDocument(lsp.DocumentURI("http://localhost:"+port), true, lsp.Range{})
	return nil
}
