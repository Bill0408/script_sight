package server

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"script_sight/controller"
)

type Server struct {
	s *http.Server
}

// handler is used to cast and function that contains
// http.ResponseWriter and *http.Request to and http.Handler.
type handler func(http.ResponseWriter, *http.Request)

// Implement ServeHTTP in order to be considered an http.Handler.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}

func New(s *http.Server) *Server {
	return &Server{
		s: s,
	}
}

func (s *Server) RegisterRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/", controller.HomePageHandler)
	r.Handle("/upload", alice.New(controller.ImageValidator, controller.ImageModifier).
		Then(handler(controller.ImageHandler)))

	s.s.Handler = r
}

func (s *Server) ListenAndServe() {
	if err := s.s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
