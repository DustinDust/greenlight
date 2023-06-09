package main

import (
	"fmt"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Creating a movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(r)
	if err != nil {
		http.NotFound(w, r)
	}
	fmt.Fprintf(w, "Showing detail of movie %d\n", id)
}
