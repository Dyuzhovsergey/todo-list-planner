package db

import (
	"database/sql"
	"modernc.org/sqlite"
	"os"
)


dbFile := "scheduler.db"
_, err := os.Stat(dbFile)

install := os.IsNotExist(err)

	// Открываем или создаём базу данных
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Если файла не было, создаём таблицу и индексы
	if install {
		createTableSQL := `
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			task_name TEXT NOT NULL,
			schedule_time DATETIME NOT NULL,
			status TEXT DEFAULT 'pending'
		);
		CREATE INDEX IF NOT EXISTS idx_schedule_time ON scheduler(schedule_time);
		`