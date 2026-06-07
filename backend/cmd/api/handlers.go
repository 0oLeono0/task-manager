package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func (app *application) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := app.postgresStore.getTasks()
	if err != nil {
		app.logger.Println(err)
		app.errorJSON(w, http.StatusInternalServerError, "server error")
		return
	}

	err = app.writeJSON(w, http.StatusOK, tasks)
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

	foundTask, err := app.postgresStore.getTaskByID(id)
	if err != nil {
		if errors.Is(err, errTaskNotFound) {
			app.errorJSON(w, http.StatusNotFound, "task not found")
			return
		}
		app.logger.Println(err)
		app.errorJSON(w, http.StatusInternalServerError, "server error")
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

	createdTask, err := app.postgresStore.createTask(inputTask.Title, inputTask.Completed)
	if err != nil {
		app.logger.Println(err)
		app.errorJSON(w, http.StatusInternalServerError, "server error")
		return
	}

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

	updatedTask, err := app.postgresStore.updateTask(id, inputTask.Title, inputTask.Completed)
	if err != nil {
		if errors.Is(err, errTaskNotFound) {
			app.errorJSON(w, http.StatusNotFound, "task not found")
			return
		}
		app.logger.Println(err)
		app.errorJSON(w, http.StatusInternalServerError, "server error")
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

	err = app.postgresStore.deleteTask(id)
	if err != nil {
		if errors.Is(err, errTaskNotFound) {
			app.errorJSON(w, http.StatusNotFound, "task not found")
			return
		}
		app.logger.Println(err)
		app.errorJSON(w, http.StatusInternalServerError, "server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
