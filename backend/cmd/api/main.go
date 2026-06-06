package main

import (
	"encoding/json"
	"errors"
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
	store  *taskStore
}

type config struct {
	serverAddr string
}

type task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type taskStore struct {
	mu     sync.Mutex
	tasks  []task
	nextID int
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
		store: &taskStore{
			tasks:  tasks,
			nextID: len(tasks) + 1,
		},
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
	r.Delete("/tasks/{id}", app.deleteTaskHandler)
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

func sampleTasks() []task {
	return []task{
		{ID: 1, Title: "Learn Go", Completed: false},
		{ID: 2, Title: "Build API", Completed: true},
	}
}

func (store *taskStore) getTasks() []task {
	store.mu.Lock()
	defer store.mu.Unlock()
	tasks := make([]task, len(store.tasks))
	copy(tasks, store.tasks)
	return tasks
}

func (store *taskStore) getTaskByID(id int) (task, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	for i := range store.tasks {
		if store.tasks[i].ID == id {
			return store.tasks[i], nil
		}
	}

	return task{}, errors.New("task not found")
}

func (store *taskStore) createTask(title string, completed bool) task {
	store.mu.Lock()
	defer store.mu.Unlock()

	task := task{
		ID:        store.nextID,
		Title:     title,
		Completed: completed,
	}

	store.tasks = append(store.tasks, task)
	store.nextID++

	return task
}

func (store *taskStore) updateTask(id int, title *string, completed *bool) (task, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	for i := range store.tasks {
		if store.tasks[i].ID == id {
			if title != nil {
				store.tasks[i].Title = *title
			}

			if completed != nil {
				store.tasks[i].Completed = *completed
			}

			return store.tasks[i], nil
		}
	}

	return task{}, errors.New("task not found")
}

func (store *taskStore) deleteTask(id int) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	for i := range store.tasks {
		if store.tasks[i].ID == id {
			store.tasks = append(store.tasks[:i], store.tasks[i+1:]...)
			return nil
		}
	}

	return errors.New("task not found")
}

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func (app *application) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := app.store.getTasks()

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

	foundTask, err := app.store.getTaskByID(id)
	if err != nil {
		app.errorJSON(w, http.StatusNotFound, "task not found")
		return
	}
	err = app.writeJSON(w, http.StatusOK, foundTask)
	if err != nil {
		app.logger.Println(err)
		return
	}
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

	createdTask := app.store.createTask(inputTask.Title, inputTask.Completed)
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

	updatedTask, err := app.store.updateTask(id, inputTask.Title, inputTask.Completed)
	if err != nil {
		app.errorJSON(w, http.StatusNotFound, "task not found")
		return
	}
	err = app.writeJSON(w, http.StatusOK, updatedTask)
	if err != nil {
		app.logger.Println(err)
		return
	}
}

func (app *application) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = app.store.deleteTask(id)
	if err != nil {
		app.errorJSON(w, http.StatusNotFound, "task not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
