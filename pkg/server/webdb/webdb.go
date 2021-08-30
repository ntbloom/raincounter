package webdb

import "database/sql"

// WebDB puts distance between database implementation and the
type WebDB interface {
	// MakeSchema populates the schema in the web database
	MakeSchema() (sql.Result, error)

	// EnterData runs arbitary sql commands
	EnterData(string) (sql.Result, error)

	// AddTagValue puts a single tag and value in the database
	AddTagValue(int, int) (sql.Result, error)

	// AddRainEvent puts a rain event with a timestamp from the sensor
	AddRainEvent(float32, string) (sql.Result, error)
}
