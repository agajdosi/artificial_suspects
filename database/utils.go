package database

import "time"

// Get Timestamp of current time.
func TimestampNow() string {
	return time.Now().Format(TimeFormat)
}
