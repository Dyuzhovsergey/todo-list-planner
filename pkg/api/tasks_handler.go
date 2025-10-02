package api

import (
	"net/http"
	"strconv"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/bl"
	"github.com/Dyuzhovsergey/todo-list-project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")
	limitStr := r.URL.Query().Get("limit")

	limit := 50

	if limitStr != "" {
		if lim, err := strconv.Atoi(limitStr); err == nil && lim > 0 && lim < 50 {
			limit = lim
		}
	}

	tasks, err := db.TasksWithSearch(limit, search)
	if err != nil {
		bl.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	bl.WriteJSON(w, TasksResp{Tasks: tasks})
}
