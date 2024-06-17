package logger

import "time"

func tz() *time.Location {
	t, _ := time.LoadLocation("Asia/Jakarta")

	return t
}
