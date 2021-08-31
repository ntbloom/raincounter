package webdb

import (
	"database/sql"
	"time"
)

// DBEntry enters data into the database
type DBEntry interface {
	// RunCmd runs arbitary sql commands
	RunCmd(string) (sql.Result, error)

	// AddTagValue puts a single tag and value in the database
	AddTagValue(int, int) (sql.Result, error)

	// AddRainEvent puts a rain event with a timestamp from the sensor
	AddRainEvent(float32, string) (sql.Result, error)
}

// DBQuery retreives data from the database
type DBQuery interface {
	// RunCmd runs arbitary sql commands
	RunCmd(string) (sql.Result, error)

	// Unwrap converts a query with one row result into the correct value
	Unwrap(sql.Result) interface{}

	// TallyRainSince gets total rain from a time in the past to present
	TallyRainSince(time.Time) float32

	// TallyRainFrom gets total rain between two timestamps
	TallyRainFrom(time.Time, time.Time) float32

	// GetLastRainTime shows the date of the last rain
	GetLastRainTime() time.Time
}
