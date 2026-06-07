package main

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func clearTasks(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec("TRUNCATE TABLE tasks RESTART IDENTITY")
	if err != nil {
		t.Fatal(err)
	}
}

func connectToDB(t *testing.T) *postgresTaskStore {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN is not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}

	store := &postgresTaskStore{
		db: db,
	}

	clearTasks(t, store.db)

	t.Cleanup(func() {
		db.Close()
	})

	return store
}

func TestPostgresTaskStore_Create(t *testing.T) {
	store := connectToDB(t)

	created, err := store.createTask("test title", false)
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 {
		t.Fatal("expected task ID to be set")
	}
	if created.Title != "test title" {
		t.Fatalf("expected title %q, got %q", "test title", created.Title)
	}
	if created.Completed != false {
		t.Fatalf("expected completed %t, got %t", false, created.Completed)
	}

	found, err := store.getTaskByID(created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if found.ID != created.ID {
		t.Fatalf("expected id %d, got %d", created.ID, found.ID)
	}
	if found.Title != created.Title {
		t.Fatalf("expected title %q, got %q", created.Title, found.Title)
	}
	if found.Completed != created.Completed {
		t.Fatalf("expected completed %v, got %v", created.Completed, found.Completed)
	}
}

func TestPostgresTaskStore_GetTasks(t *testing.T) {
	store := connectToDB(t)
	created1, err := store.createTask("test title 1", false)
	if err != nil {
		t.Fatal(err)
	}

	created2, err := store.createTask("test title 2", true)
	if err != nil {
		t.Fatal(err)
	}

	tasks, err := store.getTasks()
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].ID != created1.ID {
		t.Fatalf("expected id %d, got %d", created1.ID, tasks[0].ID)
	}
	if tasks[1].ID != created2.ID {
		t.Fatalf("expected id %d, got %d", created2.ID, tasks[1].ID)
	}

	if tasks[0].Title != created1.Title {
		t.Fatalf("expected Title %s, got %s", created1.Title, tasks[0].Title)
	}
	if tasks[1].Title != created2.Title {
		t.Fatalf("expected Title %s, got %s", created2.Title, tasks[1].Title)
	}

	if tasks[0].Completed != created1.Completed {
		t.Fatalf("expected Completed %t, got %t", created1.Completed, tasks[0].Completed)
	}
	if tasks[1].Completed != created2.Completed {
		t.Fatalf("expected Completed %t, got %t", created2.Completed, tasks[1].Completed)
	}
}
