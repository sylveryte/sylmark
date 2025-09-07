package server

import "net/http"

func (server *Server) Hello(w http.ResponseWriter, r *http.Request) {

	type Resp struct {
		Hi string `json:"hi"`
	}

	WriteJson(Resp{
		Hi: "Holamigo",
	}, w)
}
