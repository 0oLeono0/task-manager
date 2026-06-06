package main

import (
	"errors"
	"sync"
)

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

type taskStore struct {
	mu     sync.Mutex
	tasks  []task
	nextID int
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