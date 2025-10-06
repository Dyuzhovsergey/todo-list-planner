package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat)
	          VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat
	FROM scheduler
	WHERE id = ?`

	row := DB.QueryRow(query, id)
	t := &Task{}
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task %s not found: %w", id, ErrTaskNotFound)
		}
		return nil, err
	}
	return t, nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	res, err := DB.Exec(query, i)
	if err != nil {
		return fmt.Errorf("delete task failed: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("no task with id %s: %w", id, ErrTaskNotFound)
	}
	return nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler
	SET date = ?, title = ?, comment = ?, repeat = ?
	WHERE id = ?`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

// TasksWithSearch ищет задачи с учётом фильтра search и лимита.
func TasksWithSearch(limit int, search string) ([]*Task, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if search == "" {
		// просто вернуть ближайшие задачи
		query := `SELECT id, date, title, comment, repeat
		          FROM scheduler
		          ORDER BY date ASC
		          LIMIT ?`
		rows, err = DB.Query(query, limit)
	} else {
		// поиск: либо по дате (dd.mm.yyyy), либо по подстроке
		if d, errDate := time.Parse("02.01.2006", search); errDate == nil {
			searchDate := d.Format("20060102")
			query := `SELECT id, date, title, omment, repeat
			FROM scheduler
			WHERE date = ?
			LIMIT ?`
			rows, err = DB.Query(query, searchDate, limit)
		} else {
			pattern := "%" + search + "%"
			query := `SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE title LIKE ? OR comment LIKE ?
			ORDER BY date ASC
			LIMIT ?`
			rows, err = DB.Query(query, pattern, pattern, limit)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

func UpdateDate(id string, nextDate string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := DB.Exec(query, nextDate, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating date")
	}
	return nil
}
