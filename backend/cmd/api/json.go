package main

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Err string `json:"error"`
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func (app *application) errorJSON(w http.ResponseWriter, status int, message string) {
	err := app.writeJSON(w, status, errorResponse{message})
	if err != nil {
		app.logger.Println(err)
	}
}

func (app *application) readJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(dst)
}
