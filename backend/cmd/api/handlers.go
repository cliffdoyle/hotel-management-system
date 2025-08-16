// File: backend/cmd/api/handlers.go

package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// A simple healthcheck response
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
	}

	js, err := json.Marshal(data)
	if err != nil {
		app.logger.Error("failed to marshal healthcheck data", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}