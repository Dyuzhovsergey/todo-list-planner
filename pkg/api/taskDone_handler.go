package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/db"
)

func doneTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, fmt.Errorf("id task is empty"), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}

	// если задача одноразовая — удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{}) // {}
		return
	}

	// если задача повторяющаяся — вычисляем следующую дату
	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	if err := db.UpdateDate(id, next); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{}) // {}
}
