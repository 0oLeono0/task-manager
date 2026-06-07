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
