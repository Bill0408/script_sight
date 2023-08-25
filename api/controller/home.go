package controller

import (
	"fmt"
	"html/template"
	"net/http"
)

type data struct {
	Hash string
}

func (c *Controller) HomePageHandler(w http.ResponseWriter, r *http.Request) {
	d := data{Hash: c.Hash}

	// Create a new html template from the home.html file.
	t, err := template.ParseFiles("/frontend/pages/home.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "An error occurred parsing the homepage file.", http.StatusInternalServerError)
		return
	}

	// Send back the parsed home.html page.
	if err = t.Execute(w, d); err != nil {
		fmt.Println(err)
		http.Error(w, "An error occurred executing the template.", http.StatusInternalServerError)
		return
	}
}
