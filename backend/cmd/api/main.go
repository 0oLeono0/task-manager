package main

import (
	"flag"
	"log"
	"net/http"
)

type application struct {
	logger *log.Logger
	cfg    config
	store  *taskStore
}

type config struct {
	serverAddr string
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
