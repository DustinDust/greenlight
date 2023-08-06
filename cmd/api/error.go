package main

import (
	"fmt"
	"greenlight/internal/constant"
	"net/http"
)

// generic helper to logging an error
func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err.Error(), nil)
}

// generic helper to sending json formatted error to client with given status code
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}
	err := app.writeJSONResponse(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// will be used when server encountered unexpected error during runtime
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	app.errorResponse(w, r, http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusNotFound, constant.NOT_FOUND)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf(constant.METHOD_NOT_ALLOWED, r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, msg)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationError(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request) {
	message := "unable to perform action due to conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded, please try again later"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}
