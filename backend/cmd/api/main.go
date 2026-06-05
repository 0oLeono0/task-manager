package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type application struct {
	logger *log.Logger
}

func main() {
	app := &application{
		logger: log.Default(),
	}

	r := app.router()

	app.logger.Fatal(http.ListenAndServe(":8080", r))
}

func (app *application) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.healthHandler)
	return r
}

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}