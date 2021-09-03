package webdb

import (
	"time"
)

// DBEntry enters data into the database
type DBEntry interface {
	// Insert runs arbitrary sql INSERT commands
	Insert(string) error

	// AddTagValue puts a single tag and value in the database
	AddTagValue(int, int, time.Time) error

	// AddTempCValue puts a Celsius temperature value in the database
	AddTempCValue(int, time.Time) error

	// AddRainMMEvent puts a rain event with a timestamp from the sensor
	AddRainMMEvent(float32, time.Time) error

	// Close closes the connection with the database. Necessary for pooled connections
	Close()
}

// DBQuery retreives data from the database
type DBQuery interface {
	// Select runs arbitary sql SELECT commands
	Select(string) (interface{}, error)

	// TotalRainMMSince gets total rain from a time in the past to present
	TotalRainMMSince(time.Time) float32

	// TotalRainMMFrom gets total rain between two timestamps
	TotalRainMMFrom(time.Time, time.Time) float32

	// GetRainMMSince gets a RainMMMap from a time in the past to present
	GetRainMMSince(time.Time) *RainMMMap

	// GetRainMMFrom gets a RainMMMap between two timestamps
	GetRainMMFrom(time.Time, time.Time) *RainMMMap

	// GetLastRainTime shows the date of the last rain
	GetLastRainTime() time.Time

	// GetTempDataCSince gets a TempCMap from a time in the past to the present
	GetTempDataCSince(time.Time) *TempCMap

	// GetTempDataCFrom gets a TempCMap between two timestamps
	GetTempDataCFrom(time.Time, time.Time) *TempCMap

	// GetLastTempC shows the most recent temperature
	GetLastTempC() int

	// Close closes the connection with the database. Necessary for pooled connections
	Close()
}

// RainMMMap is a simple map of a timestamp and millimeters of rain
type RainMMMap map[time.Time]float32

// TempCMap is a simple map of Celsius temperatures over time
type TempCMap map[time.Time]int
