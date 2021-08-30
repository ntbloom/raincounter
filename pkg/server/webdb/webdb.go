package webdb

import (
	"database/sql"
	"time"
)

// DBEntry enters data into the database
type DBEntry interface {
	// EnterData runs arbitary sql commands
	EnterData(string) (sql.Result, error)

	// AddTagValue puts a single tag and value in the database
	AddTagValue(int, int) (sql.Result, error)

	// AddRainEvent puts a rain event with a timestamp from the sensor
	AddRainEvent(float32, string) (sql.Result, error)
}

// DBQuery queries data from the database
type DBQuery interface {
	// TallyRainSince gets total rain from a time in the past to present
	TallyRainSince(time.Time) float32

	// TallyRainFrom gets total rain between two timestamps
	TallyRainFrom(time.Time, time.Time) float32

	// GetLastRain tallies the date since the last rain
	GetLastRain() time.Time
}
