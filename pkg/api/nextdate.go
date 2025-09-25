package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

// now    — время, от которого ищется ближайшая дата
// dstart — исходное время в формате 20060102, от которого начинается отсчёт повторений
// repeat — правило повторения в описанном выше формате ("d <число>", "y", "w" - неделя)
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("reeat rule is empty")
	}

	date, err := time.Parse(DateLayout, dstart)
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

		if date.Format(DateLayout) == dstart && date.After(now) {
			date = date.AddDate(0, 0, interval)
		}

		for !date.After(now) {
			date = date.AddDate(0, 0, interval)
		}

	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid repeat formal for 'y'")
		}
		if date.After(now) {
			// если старт позже now, то нужно не его, а следующий
			date = date.AddDate(1, 0, 0)
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

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("invalid repeat format for 'm'")
		}
		var days [32]bool
		var dayLast, dayPrevLast bool

		daysStr := strings.Split(parts[1], ",")
		for _, ds := range daysStr {
			d, err := strconv.Atoi(ds)
			if err != nil {
				return "", errors.New("invalid day in month")
			}
			switch {
			case d >= 1 && d <= 31:
				days[d] = true
			case d == -1:
				dayLast = true
			case d == -2:
				dayPrevLast = true
			default:
				return "", errors.New("invalid day in month")
			}
		}

		var months [13]bool
		if len(parts) == 3 {
			monthsStr := strings.Split(parts[2], ",")
			for _, ms := range monthsStr {
				m, err := strconv.Atoi(ms)
				if err != nil || m < 1 || m > 12 {
					return "", errors.New("invalid month")
				}
				months[m] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				months[i] = true
			}
		}

		for {
			year, month, _ := date.Date()
			dim := daysInMonth(year, month)
			for d := 1; d <= dim; d++ {
				if days[d] {
					condidateDay := time.Date(year, month, d, 0, 0, 0, 0, time.UTC)
					if condidateDay.After(now) && months[int(month)] {
						return condidateDay.Format(DateLayout), nil
					}
				}
			}

			if dayLast {
				condidateDay := time.Date(year, month, dim, 0, 0, 0, 0, time.UTC)
				if condidateDay.After(now) && months[int(month)] {
					return condidateDay.Format(DateLayout), nil
				}
			}

			if dayPrevLast && dim > 1 {
				condidateDay := time.Date(year, month, dim-1, 0, 0, 0, 0, time.UTC)
				if condidateDay.After(now) && months[int(month)] {
					return condidateDay.Format(DateLayout), nil
				}
			}

			date = date.AddDate(0, 1, 0)
		}

	default:
		return "", errors.New("unsupported repeat format")
	}

	return date.Format(DateLayout), nil
}

func weekdayGoToRule(w time.Weekday) int {
	if w == time.Sunday {
		return 7
	}
	return int(w)
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
