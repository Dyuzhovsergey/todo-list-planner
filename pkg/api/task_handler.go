package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)

	default:
		http.Error(w, "metod not allowed", http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("addTaskHandler called")
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("decode error:", err)
		writeError(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("decoded task: %+v\n", task)

	if task.Title == "" {
		writeError(w, fmt.Errorf("title task is empty"), http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		log.Println("checkDate error:", err)
		writeError(w, err, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.Println("AddTask error:", err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

// GET /api/task?id=...
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, fmt.Errorf("id task is empty"), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, fmt.Errorf("task not fund"), http.StatusNotFound)
		} else {
			writeError(w, err, http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, task)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, fmt.Errorf("id task is empty"), http.StatusBadRequest)
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeError(w, err, http.StatusNotFound)
	}
	writeJSON(w, map[string]string{})
}

// PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, fmt.Errorf("error parse JSON: %w", err), http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeError(w, fmt.Errorf("id task is requared"), http.StatusBadRequest)
		return
	}

	// проверки такие же, как в addTaskHandler
	if task.Title == "" {
		writeError(w, fmt.Errorf("title task is empty"), http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, fmt.Errorf("checkDate error: %w", err), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		if strings.Contains(err.Error(), "incorrect id") {
			writeError(w, err, http.StatusNotFound) // 404
		} else {
			writeError(w, err, http.StatusInternalServerError) // 500
		}
		return
	}

	writeJSON(w, map[string]string{}) // пустой JSON {}
}

func doneTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
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

	// если repeat пустой → одноразовая задача, удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{}) // {}
		return
	}

	// если задача периодическая → считаем следующую дату
	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	// обновляем дату
	if err := db.UpdateDate(id, next); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{}) // {}
}
