package main

import "github.com/julienschmidt/httprouter"

type Server struct {
	Router *httprouter.Router
}

func NewServer() *Server {
	return &Server{
		Router: httprouter.New(),
	}
}
