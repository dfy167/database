package database

import (
	"errors"
	"time"
)

const (
	timeStringMonth  string = "200601"
	timeStringDay    string = "20060102"
	timeStringHour   string = "2006010215"
	timeStringMinute string = "200601021504"
)

type timeType uint

//  time
const (
	Month timeType = iota + 1
	Day
	Hour
	Minute
)

// String2Time String2Time
func String2Time(timeString string) (time.Time, error) {
	l := len(timeString)
	switch l {
	case 6:
		return time.Parse(timeStringMonth, timeString)
	case 8:
		return time.Parse(timeStringDay, timeString)
	case 10:
		return time.Parse(timeStringHour, timeString)
	case 12:
		return time.Parse(timeStringMinute, timeString)
	default:
		return time.Now(), errors.New("wrong time string")
	}
}

// Time2String Time2String
func Time2String(t time.Time, timeType timeType) string {
	switch timeType {
	case Month:
		return t.Format(timeStringMonth)
	case Day:
		return t.Format(timeStringDay)
	case Hour:
		return t.Format(timeStringHour)
	case Minute:
		return t.Format(timeStringMinute)
	default:
		return ""
	}
}

// TimeBetween TimeBetween
func TimeBetween(st time.Time, et time.Time, deltaType timeType) (betweenTime []time.Time) {
	if st.Equal(et) {
		return append(betweenTime, st)
	}
	var deltaTime time.Duration
	switch deltaType {
	case Day:
		deltaTime = 24 * time.Hour
	case Hour:
		deltaTime = time.Hour
	case Minute:
		deltaTime = time.Minute
	}
	for st.Before(et) {
		betweenTime = append(betweenTime, st)
		st = st.Add(deltaTime)
	}
	return
}
