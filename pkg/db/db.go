// Package db for create table
package db

import (
	"log"

	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		return err
	}
	if _, err := DB.Exec(schema); err != nil {
		return err
	}
	log.Printf("using database file: %s", dbFile)
	return nil
}
