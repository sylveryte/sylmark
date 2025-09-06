package server

import "github.com/go-chi/chi/v5"

func (server *Server) SetupRoutes(r chi.Router) {
	if server == nil {
		return
	}
	s := server
	r.Get("/hello", s.Hello)
	r.Get("/graph", s.GetGraph)
	r.Post("/document/show", s.ShowDocument)
}
