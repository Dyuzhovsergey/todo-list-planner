package main

import (
	"log"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
