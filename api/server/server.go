package server

import "net/http"

type Server struct {
	s *http.Server
}

func New(s *http.Server) *Server {
	return &Server{
		s: s,
	}
}
