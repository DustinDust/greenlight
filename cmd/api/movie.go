package main

import (
	"greenlight/internal/data"
	"greenlight/internal/validator"
	"net/http"
	"time"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	movie := data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: data.Runtime(input.Runtime),
		Genres:  input.Genres,
	}

	if data.ValidateMovieInput(v, movie); !v.Valid() {
		app.failedValidationError(w, r, v.Errors)
		return
	}

	app.writeJSONResponse(w, http.StatusOK, envelope{"movie": input}, nil)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}
	movie := data.Movie{
		ID:        id,
		Title:     "Harry Potter",
		CreatedAt: time.Now(),
		Runtime:   120,
		Genres:    []string{"drama", "fantasy", "teens"},
		Version:   1,
	}
	err = app.writeJSONResponse(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
