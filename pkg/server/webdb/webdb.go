package webdb

import (
	"database/sql"
	"time"
)

// DBEntry enters data into the database
type DBEntry interface {
	// Insert runs arbitrary sql INSERT commands
	Insert(string) error

	// AddTagValue puts a single tag and value in the database
	AddTagValue(int, int) (sql.Result, error)

	// AddRainEvent puts a rain event with a timestamp from the sensor
	AddRainEvent(float32, string) (sql.Result, error)

	// Close closes the connection with the database. This may or may not need
	// to be called, depending on the implementation. In the case of a pooled
	// connection struct, this is always necessary.
	Close()
}

// DBQuery retreives data from the database
type DBQuery interface {
	// Select runs arbitary sql SELECT commands
	Select(string) (interface{}, error)

	// TallyRainSince gets total rain from a time in the past to present
	TallyRainSince(time.Time) float32

	// TallyRainFrom gets total rain between two timestamps
	TallyRainFrom(time.Time, time.Time) float32

	// GetLastRainTime shows the date of the last rain
	GetLastRainTime() time.Time

	// Close closes the connection with the database. This may or may not need
	// to be called, depending on the implementation. In the case of a pooled
	// connection struct, this is always necessary.
	Close()
}
