package database

import "time"

const (
	TimeFormat string = time.RFC3339Nano
)

// Get Timestamp of current time.
func TimestampNow() string {
	return time.Now().Format(TimeFormat)
}
