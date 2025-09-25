package main

import (
	"log"
	"os"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/db"
	"github.com/Dyuzhovsergey/todo-list-project/pkg/server"
)

const (
	defaultDBFile = "data/scheduler.db"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("cannot init database: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("cannot run server: %v", err)
	}

}
