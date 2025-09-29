package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("addTaskHandler called")
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("decode error:", err)
		writeJSON(w, map[string]string{"error": "error parse JSON: " + err.Error()})
		return
	}

	log.Printf("decoded task: %+v\n", task)

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "title task is empty"})
		return
	}

	if err := checkDate(&task); err != nil {
		log.Println("checkDate error:", err)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		log.Println("AddTask error:", err)
		writeJSON(w, map[string]string{"error": "error add task" + err.Error()})
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(data)
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
