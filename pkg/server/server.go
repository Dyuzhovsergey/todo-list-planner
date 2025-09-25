package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/api"
)

const (
	webDir          = "web"
	defaultHttpPort = "7540"
)

func Run() error {
	httpPort := os.Getenv("TODO_PORT")
	if httpPort == "" {
		httpPort = defaultHttpPort
	}

	api.Init()

	fileHandler := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileHandler)

	addr := ":" + httpPort
	fmt.Printf("Server TODO started at http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}
