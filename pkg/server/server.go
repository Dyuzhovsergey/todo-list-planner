package server

import (
	"fmt"
	"net/http"
	"os"
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
	addr := ":" + httpPort

	fileHandler := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileHandler)

	fmt.Printf("Server TODO started at http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}
