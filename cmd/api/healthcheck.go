package main

import (
	"net/http"
	"time"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": envelope{
			"environment": app.config.env,
			"version":     version,
		},
	}

	time.Sleep(4 * time.Second)

	err := app.writeJSONResponse(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
