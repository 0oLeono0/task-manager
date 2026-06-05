package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type application struct {
	logger *log.Logger
	cfg config
}

type config struct {
	serverAddr string
}

func main() {
	cfg := config{
		serverAddr: ":8080",
	}

	app := &application{
		logger: log.Default(),
		cfg: cfg,
	}

	r := app.router()

	app.logger.Fatal(http.ListenAndServe(app.cfg.serverAddr, r))
}

func (app *application) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.healthHandler)
	return r
}

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}