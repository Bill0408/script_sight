package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"script_sight/controller"
	"script_sight/server"
)

const (
	port = ":8080"

	dirName  = "../script_sight/api/img"
	ownerRWS = 0700
)

func main() {
	createImgDir()

	c := controller.New(hashFile("../script_sight/frontend/static/home.js"))

	s := &http.Server{Addr: port}

	// Initialize the server with configuration.
	srv := server.New(s, c)

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

// hashFile takes in a filepath and creates a hash of the file's content.
func hashFile(path string) string {
	file, err := os.Open(path) // Open the file at the specified path.
	if err != nil {
		log.Fatal(err)
	}
	// Release resources when hashfile returns.
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			os.Exit(1)
		}
	}(file)

	hash := sha256.New() // Get a new hash computing the SHA256 checksum.
	// Copy the content from the file into the hash.
	if _, err = io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}

	hashInBytes := hash.Sum(nil)[:20]            // Get the hash.
	hashString := fmt.Sprintf("%x", hashInBytes) // Convert the hash to string.

	return hashString[0:20]
}
