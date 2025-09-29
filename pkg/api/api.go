// Package api provides a wonderful application for doing amazing things.
package api

import (
	"net/http"
)

func InitHandlers() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)

}
