package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type application struct {
	logger *log.Logger
	cfg    config
}

type config struct {
	serverAddr string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.serverAddr, "addr", ":8080", "HTTP server address")

	flag.Parse()

	app := &application{
		logger: log.Default(),
		cfg:    cfg,
	}

	r := app.router()

	app.logger.Printf("server starting on address: %s", app.cfg.serverAddr)
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
