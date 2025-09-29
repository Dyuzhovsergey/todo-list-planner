// Package server for create server work
package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/api"
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
	httpPort := os.Getenv("TODO_PORT")
	if httpPort == "" {
		httpPort = defaultHTTPPort
	}

	api.InitHandlers()
	registerStatic()

	addr := ":" + httpPort
	fmt.Printf("Server TODO started at http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}
