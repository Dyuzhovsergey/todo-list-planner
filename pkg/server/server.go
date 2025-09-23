package server

import (
	"fmt"
	"net/http"
)

const (
	webDir   = "web"
	httpPort = ":7540"
)

func Run() error {
	fileHandler := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileHandler)

	fmt.Printf("Server started at http://localhost%s\n", httpPort)
	return http.ListenAndServe(httpPort, nil)
}
