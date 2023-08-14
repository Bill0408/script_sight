package server

import (
	"github.com/gorilla/mux"
	"log"
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

func (s *Server) ListenAndServe() {
	if err := s.s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
