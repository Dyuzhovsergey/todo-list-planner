package main

import (
	"log"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/db"
	"github.com/Dyuzhovsergey/todo-list-project/pkg/server"
)

const (
	dbFile = "scheduler.db"
)

func main() {

	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("cannot init database: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("cannot run server: %v", err)
	}

}
