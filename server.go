package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Router *httprouter.Router
}

func NewServer() *Server {
	return &Server{
		Router: httprouter.New(),
	}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
