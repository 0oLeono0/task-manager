package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type fakeStore struct {
	tasks []task
}

func (p *fakeStore) getTasks() ([]task, error) {
	return p.tasks, nil
}

func (p *fakeStore) getTaskByID(id int) (task, error) {
	for _, task := range p.tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return task{}, errTaskNotFound
}

func (p *fakeStore) createTask(title string, completed bool) (task, error) {
	return task{}, nil
}

func (p *fakeStore) updateTask(id int, title *string, completed *bool) (task, error) {
	return task{}, nil
}

func (p *fakeStore) deleteTask(id int) error {
	return nil
}

func TestHealthHandler(t *testing.T) {
	app := &application{}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Body.String() != "ok\n" {
		t.Fatalf("expected body %q, got %q", "ok\n", rec.Body.String())
	}
}

func TestGetTasks(t *testing.T) {
	testTasks := []task{
		{ID: 1, Title: "test title 1", Completed: false},
		{ID: 2, Title: "test title 2", Completed: true},
	}

	app := &application{
		taskStore: &fakeStore{
			tasks: testTasks,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected %s, got %s", "application/json", contentType)
	}

	var body []task
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(body, testTasks) {
		t.Fatalf("expected %+v, got %+v", testTasks, body)
	}
}

func TestGetTaskByID(t *testing.T) {
	tsk := task{ID: 1, Title: "test title 1", Completed: false}
	app := &application{
		taskStore: &fakeStore{
			tasks: []task{tsk},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected %s, got %s", "application/json", contentType)
	}

	var body task
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(body, tsk) {
		t.Fatalf("expected %+v, got %+v", tsk, body)
	}
}

func TestGetTaskByID_NotFound(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{},
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, rec.Code)
	}
}
