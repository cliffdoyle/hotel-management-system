// File: backend/cmd/api/handlers.go

package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// refactor to use the new envelope type
	data := envelope{
		"status":      "available",
		"environment": app.config.env,
	}

	//use the writeJSON helper instead now

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		// app.logger.Error("failed to marshal healthcheck data", "error", err)
		// http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		// return
		//Use the serverErrorResponse helper instead
		app.serverErrorResponse(w, r, err)
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.Write(js)
}
