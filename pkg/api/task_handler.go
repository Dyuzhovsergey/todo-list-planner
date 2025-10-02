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

	"github.com/Dyuzhovsergey/todo-list-project/pkg/bl"
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
		bl.WriteError(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("decoded task: %+v\n", task)

	if task.Title == "" {
		bl.WriteError(w, fmt.Errorf("title task is empty"), http.StatusBadRequest)
		return
	}

	if err := bl.CheckDate(&task); err != nil {
		log.Println("checkDate error:", err)
		bl.WriteError(w, err, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.Println("AddTask error:", err)
		bl.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	bl.WriteJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		bl.WriteError(w, fmt.Errorf("id task is empty"), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			bl.WriteError(w, fmt.Errorf("task not fund"), http.StatusNotFound)
		} else {
			bl.WriteError(w, err, http.StatusInternalServerError)
		}
		return
	}

	bl.WriteJSON(w, task)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		bl.WriteError(w, fmt.Errorf("id task is empty"), http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		bl.WriteError(w, err, http.StatusNotFound)
		return
	}

	bl.WriteJSON(w, map[string]string{})
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		bl.WriteError(w, fmt.Errorf("error parse JSON: %w", err), http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		bl.WriteError(w, fmt.Errorf("id task is requared"), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		bl.WriteError(w, fmt.Errorf("title task is empty"), http.StatusBadRequest)
		return
	}

	if err := bl.CheckDate(&task); err != nil {
		bl.WriteError(w, fmt.Errorf("checkDate error: %w", err), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		if strings.Contains(err.Error(), "incorrect id") {
			bl.WriteError(w, err, http.StatusNotFound)
		} else {
			bl.WriteError(w, err, http.StatusInternalServerError)
		}
		return
	}

	bl.WriteJSON(w, map[string]string{})
}
