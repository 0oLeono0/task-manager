package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	logger        *log.Logger
	cfg           config
	memoryStore   *taskStore
	postgresStore *postgresTaskStore
}

type config struct {
	serverAddr string
	dsn        string
}

func main() {
	var cfg config
	logger := log.Default()
	flag.StringVar(&cfg.serverAddr, "addr", ":8080", "HTTP server address")
	flag.StringVar(&cfg.dsn, "db-dsn", "", "DSN")

	flag.Parse()
	if cfg.dsn == "" {
		logger.Fatal("dsn is empty")
	}

	db, err := sql.Open("pgx", cfg.dsn)
	if err != nil {
		logger.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Print("success connect to db")

	tasks := sampleTasks()
	app := &application{
		logger: logger,
		cfg:    cfg,
		memoryStore: &taskStore{
			tasks:  tasks,
			nextID: len(tasks) + 1,
		},
		postgresStore: &postgresTaskStore{
			db: db,
		},
	}

	r := app.router()

	app.logger.Printf("server starting on address: %s", app.cfg.serverAddr)
	app.logger.Fatal(http.ListenAndServe(app.cfg.serverAddr, r))
}
