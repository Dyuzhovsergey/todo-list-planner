package main

import (
	"log"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/config"
	"github.com/Dyuzhovsergey/todo-list-project/pkg/db"
	"github.com/Dyuzhovsergey/todo-list-project/pkg/server"
)

func main() {

	config.Init()

	if err := db.Init(config.Cfg.DBFile); err != nil {
		log.Fatalf("cannot init database: %v", err)
	}
	defer db.DB.Close()

	if err := server.Run(); err != nil {
		log.Fatalf("cannot run server: %v", err)
	}
}
