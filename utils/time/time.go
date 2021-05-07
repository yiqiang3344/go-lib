package cTime

import (
	"time"
)

func Datetime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func DatetimeByTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

func MilliDatetimeByTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05.000000")
}
