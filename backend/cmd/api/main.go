package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

type errorResponse struct {
	Err string `json:"error"`
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
	r.Get("/tasks/{id}", app.getTaskHandler)
	return r
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

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func sampleTasks() []task {
	return []task{
		{ID: 1, Title: "Learn Go", Completed: false},
		{ID: 2, Title: "Build API", Completed: true},
	}
}

func (app *application) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := sampleTasks()

	err := app.writeJSON(w, http.StatusOK, tasks)
	if err != nil {
		app.logger.Println(err)
		return
	}
}

func (app *application) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, "invalid task id")
		return
	}

	tasks := sampleTasks()
	for _, task := range tasks {
		if task.ID == id {
			err = app.writeJSON(w, http.StatusOK, task)
			if err != nil {
				app.logger.Println(err)
				return
			}
			return
		}
	}
	app.errorJSON(w, http.StatusNotFound, "task not found")
}
