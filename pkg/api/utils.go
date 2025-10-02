package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/db"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(data)
}

func writeError(w http.ResponseWriter, err error, status int) {
	log.Println("error:", err)
	w.WriteHeader(status)
	writeJSON(w, map[string]string{"error": err.Error()})
}

func checkDate(task *db.Task) error {
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
