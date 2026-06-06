package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

type application struct {
	logger *log.Logger
	cfg    config

	mu     sync.Mutex
	tasks  []task
	nextID int
}

type config struct {
	serverAddr string
}

type task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type createTaskInput struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type updateTaskInput struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

type errorResponse struct {
	Err string `json:"error"`
}

func main() {
	var cfg config
	flag.StringVar(&cfg.serverAddr, "addr", ":8080", "HTTP server address")

	flag.Parse()

	tasks := sampleTasks()
	app := &application{
		logger: log.Default(),
		cfg:    cfg,
		tasks:  tasks,
		nextID: len(tasks) + 1,
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
	r.Post("/tasks", app.createTaskHandler)
	r.Patch("/tasks/{id}", app.updateTaskHandler)
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

func (app *application) readJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
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
	app.mu.Lock()
	tasks := make([]task, len(app.tasks))
	copy(tasks, app.tasks)
	app.mu.Unlock()

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

	var foundTask task
	var found bool

	app.mu.Lock()
	tasks := make([]task, len(app.tasks))
	copy(tasks, app.tasks)
	for _, task := range tasks {
		if task.ID == id {
			foundTask = task
			found = true
		}
	}
	app.mu.Unlock()

	if found {
		err = app.writeJSON(w, http.StatusOK, foundTask)
		if err != nil {
			app.logger.Println(err)
			return
		}
		return
	}
	app.errorJSON(w, http.StatusNotFound, "task not found")
}

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	inputTask := createTaskInput{}
	err := app.readJSON(r, &inputTask)
	if err != nil {
		app.logger.Println(err)
		app.errorJSON(w, http.StatusBadRequest, "invalid task")
		return
	}
	if inputTask.Title == "" {
		app.errorJSON(w, http.StatusBadRequest, "title is required")
		return
	}

	app.mu.Lock()
	createdTask := task{
		ID:        app.nextID,
		Title:     inputTask.Title,
		Completed: inputTask.Completed,
	}

	app.tasks = append(app.tasks, createdTask)
	app.nextID++
	app.mu.Unlock()

	err = app.writeJSON(w, http.StatusCreated, createdTask)
	if err != nil {
		app.logger.Println(err)
		return
	}
}

func (app *application) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, "invalid id")
		return
	}

	inputTask := updateTaskInput{}
	err = app.readJSON(r, &inputTask)
	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, "invalid request")
		return
	}

	var updatedTask task
	var found bool

	app.mu.Lock()
	for i := range app.tasks {
		if app.tasks[i].ID == id {
			found = true
			if inputTask.Title != nil {
				app.tasks[i].Title = *inputTask.Title
			}

			if inputTask.Completed != nil {
				app.tasks[i].Completed = *inputTask.Completed
			}
			updatedTask = app.tasks[i]
			break
		}
	}
	app.mu.Unlock()

	if found {
		err = app.writeJSON(w, http.StatusOK, updatedTask)
		if err != nil {
			app.logger.Println(err)
		}
		return
	}
	app.errorJSON(w, http.StatusNotFound, "task not found")
}
