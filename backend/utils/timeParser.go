package utils

import "time"

func ParseTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
