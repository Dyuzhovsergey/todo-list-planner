package bl

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Dyuzhovsergey/todo-list-planner/pkg/db"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(data); err != nil {
		log.Printf("WriteJSON error: %v", err)
	}
}

func WriteJSONSuccess(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, data)
}

func WriteJSONError(w http.ResponseWriter, err error, code int) {
	writeJSON(w, code, map[string]string{"error": err.Error()})
}

func CheckDate(task *db.Task) error {
	now := time.Now()
	todayTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	today := todayTime.Format(DateLayout)

	if task.Date == "" {
		task.Date = today
	}

	parsedDate, err := time.Parse(DateLayout, task.Date)
	if err != nil {
		return fmt.Errorf(" incorrect format date: %w", err)
	}

	if task.Repeat != "" {
		nextdate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("incorrect rule repeat: %w", err)
		}
		if parsedDate.Before(todayTime) {
			task.Date = nextdate
		}

	} else {
		if parsedDate.Before(todayTime) {
			task.Date = today
		}
	}

	return nil
}
