package main

import (
	"net/http"
	"script_sight/server"
)

const port = ":8080"

func main() {
	s := &http.Server{Addr: port}

	// Initialize the server with configuration.
	srv := server.New(s)

	// Ensure routes are registered before starting the server.
	srv.RegisterRoutes()

	// Start the server
	srv.ListenAndServe()
}
