package api

import (
	"net/http"
	"time"

	"github.com/Dyuzhovsergey/todo-list-project/pkg/bl"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(bl.DateLayout, nowStr)
		if err != nil {
			http.Error(w, "invalid now format, expected YYYYMMDD", http.StatusBadRequest)
			return
		}
	}

	next, err := bl.NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
