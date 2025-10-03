// This file provides a utility function for parsing time.
package utils

// "time" provides functions for working with time. It is used here to format a time.Time object.
import "time"

// ParseTime formats a time.Time object into an RFC3339 string.
// RFC3339 is a standard format for representing dates and times.
//
// @param t time.Time - The time.Time object to be formatted.
// @return string - The formatted time string.
func ParseTime(t time.Time) string {
	// t.Format() returns a textual representation of the time value formatted according to the layout defined by the argument.
	// time.RFC3339 is a predefined layout for RFC3339 format.
	return t.Format(time.RFC3339)
}