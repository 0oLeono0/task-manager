package main

import "errors"

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

type taskStore interface {
	getTasks() ([]task, error)
	getTaskByID(id int) (task, error)
	createTask(title string, completed bool) (task, error)
	updateTask(id int, title *string, completed *bool) (task, error)
	deleteTask(id int) error
}

var errTaskNotFound = errors.New("task not found")
