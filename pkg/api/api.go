package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// now    — время, от которого ищется ближайшая дата
// dstart — исходное время в формате 20060102, от которого начинается отсчёт повторений
// repeat — правило повторения в описанном выше формате ("d <число>", "y", "w" - неделя)
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("reeat rule is empty")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", errors.New("erroe parse dstart, invalid start date format")
	}
	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid repeat format  for 'd'")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", errors.New("invalid day interval")
		}

		for !date.After(now) {
			date = date.AddDate(0, 0, interval)
		}
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid repeat formal for 'y'")
		}
		for !date.After(now) {
			date = date.AddDate(1, 0, 0)
		}
	case "w":
		if len(parts) != 2 {
			return "", errors.New("invalid repaet format for 'w'")
		}
		var week [8]bool

		days := strings.Split(parts[1], ",")
		for _, day := range days {
			dayNum, err := strconv.Atoi(day)
			if err != nil || dayNum < 1 || dayNum > 7 {
				return "", errors.New("invalid week number")
			}
			week[dayNum] = true
		}
		for !date.After(now) || week[weekdayGoToRule(date.Weekday())] {
			date = date.AddDate(0, 0, 1)
		}

	default:
		return "", errors.New("nsupported repeat format")
	}

	return date.Format("20060102"), nil
}

func weekdayGoToRule(w time.Weekday) int {
	if w == time.Sunday {
		return 7
	}
	return int(w)
}
