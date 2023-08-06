package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// handle error cases
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)

	// Fully update, replacing the whole object in database
	router.HandlerFunc(http.MethodPut, "/v1/movies/:id", app.updateMovieHandler)
	// Can take partial updates body => replacing only meaningful part
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.partialUpdateMovieHandler)

	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.findMovieHandler)

	// wrap the router with panic recovery middleware
	return app.recoverPanic(app.rateLimit(router))
}
