package server

import "net/http"

type Server struct {
	s *http.Server
}
