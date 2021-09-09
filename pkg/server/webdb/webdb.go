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

	// AddStatusUpdate adds a status message for an asset with an integer ID
	AddStatusUpdate(int, time.Time) error

	// AddRainMMEvent puts a rain event with a timestamp from the sensor
	AddRainMMEvent(float64, time.Time) error

	// Close closes the connection with the database. Necessary for pooled connections
	Close()
}

// DBQuery retreives data from the database
type DBQuery interface {
	// Select runs arbitary sql SELECT commands
	Select(string) (interface{}, error)

	// TotalRainMMSince gets total rain from a time in the past to present
	TotalRainMMSince(time.Time) (float64, error)

	// TotalRainMMFrom gets total rain between two timestamps
	TotalRainMMFrom(time.Time, time.Time) (float64, error)

	// GetRainMMSince gets a RainEntriesMm from a time in the past to present
	GetRainMMSince(time.Time) (*RainEntriesMm, error)

	// GetRainMMFrom gets a RainEntriesMm between two timestamps
	GetRainMMFrom(time.Time, time.Time) (*RainEntriesMm, error)

	// GetLastRainTime shows the date of the last rain
	GetLastRainTime() (time.Time, error)

	// GetTempDataCSince gets a TempEntriesC from a time in the past to the present
	GetTempDataCSince(time.Time) (*TempEntriesC, error)

	// GetTempDataCFrom gets a TempEntriesC between two timestamps
	GetTempDataCFrom(time.Time, time.Time) (*TempEntriesC, error)

	// GetLastTempC shows the most recent temperature
	GetLastTempC() (int, error)

	// IsGatewayUp tells whether the gateway has published a status message in a certain time
	IsGatewayUp(time.Duration) (bool, error)

	// IsSensorUp tells whether the sensor has published a status message in a certain time
	IsSensorUp(time.Duration) (bool, error)

	// Close closes the connection with the database. Necessary for pooled connections
	Close()
}

// RainEntriesMm is a simple array of RainEntryMm values
type RainEntriesMm []RainEntryMm

// RainEntryMm is a single timestamp/mm of rain entry
type RainEntryMm struct {
	Timestamp   time.Time
	Millimeters float64
}

// TempEntriesC is an ordered slice of TempEntryC values
type TempEntriesC []TempEntryC

// TempEntryC is a single temperature/timestamp entry
type TempEntryC struct {
	Timestamp time.Time
	TempC     int
}
