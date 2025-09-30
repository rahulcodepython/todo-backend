package utils

import "time"

// ParseTime formats a given time.Time object into a string representation
// using the RFC3339 format. This format is a common standard for representing
// dates and times in internet protocols and is often used in APIs for consistency.
//
// Parameters:
// - t: The time.Time object to be formatted.
//
// Returns:
// - A string representing the formatted time in RFC3339 layout.
// RFC3339 = "2006-01-02T15:04:05Z07:00"
func ParseTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
