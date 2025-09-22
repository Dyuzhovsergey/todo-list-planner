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
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	fmt.Printf("Server started at http://localhost%s\n", httpPort)
	return http.ListenAndServe(httpPort, nil)
}
