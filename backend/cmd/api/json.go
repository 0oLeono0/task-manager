package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type errorResponse struct {
	Err string `json:"error"`
}

const maxJSONBodySize = 1 << 20

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

func isRequestBodyTooLarge(err error) bool {
	var maxBytesError *http.MaxBytesError
	return errors.As(err, &maxBytesError)
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {
		return err
	}

	var extra json.RawMessage
	err = decoder.Decode(&extra)
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err != nil {
		return err
	}

	return errors.New("body must contain only one JSON value")
}
