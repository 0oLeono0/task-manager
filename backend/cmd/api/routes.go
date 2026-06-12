package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/health", app.healthHandler)
	r.Get("/tasks", app.getTasksHandler)
	r.Get("/tasks/{id}", app.getTaskHandler)
	r.Post("/tasks", app.createTaskHandler)
	r.Patch("/tasks/{id}", app.updateTaskHandler)
	r.Delete("/tasks/{id}", app.deleteTaskHandler)
	return r
}
