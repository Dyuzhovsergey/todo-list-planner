package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/bl"
	"github.com/Dyuzhovsergey/todo-list-project/pkg/db"
)

func doneTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		bl.WriteError(w, fmt.Errorf("id task is empty"), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		bl.WriteError(w, err, http.StatusNotFound)
		return
	}

	// если задача одноразовая — удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			bl.WriteError(w, err, http.StatusInternalServerError)
			return
		}
		bl.WriteJSON(w, map[string]string{}) // {}
		return
	}

	// если задача повторяющаяся — вычисляем следующую дату
	next, err := bl.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		bl.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	if err := db.UpdateDate(id, next); err != nil {
		bl.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	bl.WriteJSON(w, map[string]string{}) // {}
}
