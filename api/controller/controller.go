package controller

import (
	"html/template"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new html template from the home.html file.
	t, err := template.ParseFiles("../script_sight/frontend/pages/home.html")
	if err != nil {
		http.Error(w, "An error occurred parsing the homepage file.", http.StatusInternalServerError)
		return
	}

	// Send back the parsed home.html page.
	if err = t.Execute(w, nil); err != nil {
		http.Error(w, "An error occurred executing the template.", http.StatusInternalServerError)
		return
	}
}
