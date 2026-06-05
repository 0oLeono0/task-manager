package main

import (
	"encoding/json"
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

type task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
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
	r.Get("/tasks", app.getTasksHandler)
	return r
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func (app *application) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := []task{
		{ID: 1, Title: "Learn Go", Completed: false},
		{ID: 2, Title: "Build API", Completed: true},
	}

	err := app.writeJSON(w, http.StatusOK, tasks)
	if err != nil {
		app.logger.Println(err)
		return
	}
}
