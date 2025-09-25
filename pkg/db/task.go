package db

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	TITLE   string `json:"title"`
	COMMENT string `json:"comment"`
	REPEAT  string `json:"repeat"`
}
