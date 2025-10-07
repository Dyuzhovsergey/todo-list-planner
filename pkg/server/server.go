// Package server for work logic server
package server

import (
	"fmt"
	"net/http"

	"github.com/Dyuzhovsergey/todo-list-planner/pkg/api"
	"github.com/Dyuzhovsergey/todo-list-planner/pkg/config"
)

const webDir = "web"

func registerStatic() {
	fileHandler := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileHandler)
}

func Run() error {
	api.InitHandlers()
	registerStatic()

	addr := ":" + config.Cfg.HTTPPort
	fmt.Printf("Server TODO started at http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}
