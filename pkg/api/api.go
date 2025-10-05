// Package api provides a wonderful application for doing amazing things.
package api

import (
	"net/http"

	"github.com/Dyuzhovsergey/todo-list-planner/pkg/auth"
)

func InitHandlers() {
	http.HandleFunc("/api/signin", auth.SigninHandler)

	http.HandleFunc("/api/nextdate", auth.AuthMiddleware(nextDateHandler))
	http.HandleFunc("/api/task", auth.AuthMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", auth.AuthMiddleware(getTasksHandler))
	http.HandleFunc("/api/task/done", auth.AuthMiddleware(doneTasksHandler))
}
