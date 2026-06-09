package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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
	task := task{ID: len(p.tasks) + 1, Title: title, Completed: completed}
	p.tasks = append(p.tasks, task)
	return task, nil
}

func (p *fakeStore) updateTask(id int, title *string, completed *bool) (task, error) {
	for i, task := range p.tasks {
		if task.ID == id {
			if title != nil {
				p.tasks[i].Title = *title
			}
			if completed != nil {
				p.tasks[i].Completed = *completed
			}
			return p.tasks[i], nil
		}
	}
	return task{}, errTaskNotFound
}

func (p *fakeStore) deleteTask(id int) error {
	for i := range p.tasks {
		if p.tasks[i].ID == id {
			p.tasks = append(p.tasks[:i], p.tasks[i+1:]...)
			return nil
		}
	}
	return errTaskNotFound
}

func checkContentType(t *testing.T, r *httptest.ResponseRecorder) {
	t.Helper()

	contentType := r.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected %s, got %s", "application/json", contentType)
	}
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

	checkContentType(t, rec)

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

	checkContentType(t, rec)

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

func TestCreateTask(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{},
	}

	reqBody := strings.NewReader(`{"title": "test title", "completed": false}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", reqBody)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, rec.Code)
	}

	checkContentType(t, rec)

	var body task
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatal(err)
	}

	testTask := task{ID: 1, Title: "test title", Completed: false}
	if !reflect.DeepEqual(body, testTask) {
		t.Fatalf("expected %+v, got %+v", testTask, body)
	}
}

func TestCreateTask_TrimsTitle(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{},
	}

	reqBody := strings.NewReader(`{"title": "  test title  ", "completed": false}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", reqBody)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, rec.Code)
	}

	var body task
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatal(err)
	}

	if body.Title != "test title" {
		t.Fatalf("expected %s, got %s", "test title", body.Title)
	}
}

func TestCreateTask_InvalidPayload(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{},
		logger:    log.New(io.Discard, "", 0),
	}

	reqBody := strings.NewReader(`"error": "test title", "completed": 67`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", reqBody)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{},
		logger:    log.New(io.Discard, "", 0),
	}

	reqBody := strings.NewReader(`{"completed": false}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", reqBody)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateTask_BlankTitle(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{},
		logger:    log.New(io.Discard, "", 0),
	}

	reqBody := strings.NewReader(`{"title": "   ", "completed": false}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", reqBody)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateTask(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{
			tasks: []task{{ID: 1, Title: "test title", Completed: false}},
		},
	}

	reqBody := strings.NewReader(`{"title": "new title"}`)
	req := httptest.NewRequest(http.MethodPatch, "/tasks/1", reqBody)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}

	checkContentType(t, rec)

	var body task
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatal(err)
	}

	if body.Title != "new title" {
		t.Fatalf("expected %s, got %s", "new title", body.Title)
	}
	if body.Completed != false {
		t.Fatalf("expected %t, got %t", false, body.Completed)
	}
}

func TestUpdateTask_CompletedOnly(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{
			tasks: []task{{ID: 1, Title: "old title", Completed: false}},
		},
	}

	reqBody := strings.NewReader(`{"completed": true}`)
	req := httptest.NewRequest(http.MethodPatch, "/tasks/1", reqBody)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}

	var body task
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatal(err)
	}

	if body.Title != "old title" {
		t.Fatalf("expected %s, got %s", "old title", body.Title)
	}
	if body.Completed != true {
		t.Fatalf("expected %t, got %t", true, body.Completed)
	}
}

func TestUpdateTask_TrimsTitle(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{
			tasks: []task{{ID: 1, Title: "test title", Completed: false}},
		},
	}

	reqBody := strings.NewReader(`{"title": "  new title  "}`)
	req := httptest.NewRequest(http.MethodPatch, "/tasks/1", reqBody)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}

	var body task
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatal(err)
	}

	if body.Title != "new title" {
		t.Fatalf("expected %s, got %s", "new title", body.Title)
	}
}

func TestUpdateTask_InvalidURL(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{},
		logger:    log.New(io.Discard, "", 0),
	}

	reqBody := strings.NewReader(`{"title": "new title"}`)
	req := httptest.NewRequest(http.MethodPatch, "/tasks/abc", reqBody)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateTask_InvalidPayload(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{
			tasks: []task{{ID: 1, Title: "test title", Completed: false}},
		},
		logger: log.New(io.Discard, "", 0),
	}

	reqBody := strings.NewReader(`{error 67}`)
	req := httptest.NewRequest(http.MethodPatch, "/tasks/1", reqBody)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateTask_BlankTitle(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{
			tasks: []task{{ID: 1, Title: "test title", Completed: false}},
		},
		logger: log.New(io.Discard, "", 0),
	}

	reqBody := strings.NewReader(`{"title": "   "}`)
	req := httptest.NewRequest(http.MethodPatch, "/tasks/1", reqBody)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateTask_NotFound(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{},
		logger:    log.New(io.Discard, "", 0),
	}

	reqBody := strings.NewReader(`{"title": "new title"}`)
	req := httptest.NewRequest(http.MethodPatch, "/tasks/999", reqBody)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestDeleteTask(t *testing.T) {
	testStore := &fakeStore{
		tasks: []task{
			{ID: 1, Title: "test title", Completed: false},
		},
	}
	app := &application{
		taskStore: testStore,
	}

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d", http.StatusNoContent, rec.Code)
	}

	if rec.Body.String() != "" {
		t.Fatalf("expected %s, got %s", "", rec.Body.String())
	}

	if len(testStore.tasks) != 0 {
		t.Fatalf("expected %d, got %d", 0, len(testStore.tasks))
	}
}

func TestDeleteTask_InvalidURL(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{},
	}

	req := httptest.NewRequest(http.MethodDelete, "/tasks/abc", nil)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	app := &application{
		taskStore: &fakeStore{},
	}

	req := httptest.NewRequest(http.MethodDelete, "/tasks/999", nil)
	rec := httptest.NewRecorder()

	app.router().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, rec.Code)
	}
}
