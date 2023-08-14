package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"script_sight/controller"
)

type Server struct {
	s *http.Server
}

func New(s *http.Server) *Server {
	return &Server{
		s: s,
	}
}

func (s *Server) RegisterRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/", controller.HomePageHandler)

	s.s.Handler = r
}
