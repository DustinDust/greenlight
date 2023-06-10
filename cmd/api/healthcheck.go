package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "Available",
		"system_info": envelope{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSONResponse(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
