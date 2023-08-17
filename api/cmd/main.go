package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"script_sight/server"
)

const (
	port = ":8080"

	dirName  = "../script_sight/api/img"
	ownerRWS = 0700
)

func main() {
	createImgDir()

	s := &http.Server{Addr: port}

	// Initialize the server with configuration.
	srv := server.New(s)

	// Ensure routes are registered before starting the server.
	srv.RegisterRoutes()

	// Start the server
	srv.ListenAndServe()
}

// createImgDir creates the img directory in the controller directory
// that will be used to temporarily store images uploaded to the server.
func createImgDir() {
	err := os.Mkdir(dirName, ownerRWS)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal(err)
		}
	}
}
