package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dyuzhovsergey/todo-list-planner/pkg/bl"
	"github.com/Dyuzhovsergey/todo-list-planner/pkg/db"
)

func doneTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		bl.WriteJSONError(w, fmt.Errorf("id task is empty"), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		bl.WriteJSONError(w, err, http.StatusNotFound)
		return
	}

	// если задача одноразовая — удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			bl.WriteJSONError(w, err, http.StatusInternalServerError)
			return
		}
		bl.WriteJSONSuccess(w, map[string]string{}) // {}
		return
	}

	// если задача повторяющаяся — вычисляем следующую дату
	next, err := bl.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		bl.WriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	if err := db.UpdateDate(id, next); err != nil {
		bl.WriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	bl.WriteJSONSuccess(w, map[string]string{}) // {}
}
