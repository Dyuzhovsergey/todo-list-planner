// Package server for work logic server
package server

import (
	"fmt"
	"net/http"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/api"
	"github.com/Dyuzhovsergey/todo-list-project/pkg/config"
)

const (
	webDir          = "web"
	defaultHTTPPort = "7540"
)

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
