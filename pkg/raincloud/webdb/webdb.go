package webdb

import (
	"time"
)

// DBEntry enters data into the database
type DBEntry interface {
	// Insert runs arbitrary sql INSERT commands
	Insert(cmd string) error

	// AddTagValue puts a single tag and value in the database
	AddTagValue(tag int, value int, gwTimestamp time.Time) error

	// AddTempCValue puts a Celsius temperature value in the database
	AddTempCValue(tempC int, gwTimestamp time.Time) error

	// AddStatusUpdate adds a status message for an asset with an integer ID
	AddStatusUpdate(asset int, gwTimeamp time.Time) error

	// AddRainMMEvent puts a rain event with a timestamp from the sensor
	AddRainMMEvent(amount float64, gwTimestamp time.Time) error

	// Close closes the connection with the database. Necessary for pooled connections
	Close()
}

// DBQuery retreives data from the database
type DBQuery interface {
	// Select runs arbitary sql SELECT commands
	Select(cmd string) (interface{}, error)

	// TotalRainMMSince gets total rain from a time in the past to present
	TotalRainMMSince(since time.Time) (float64, error)

	// TotalRainMMFrom gets total rain between two timestamps
	TotalRainMMFrom(from time.Time, to time.Time) (float64, error)

	// GetRainMMSince gets a RainEntriesMm from a time in the past to present
	GetRainMMSince(since time.Time) (*RainEntriesMm, error)

	// GetRainMMFrom gets a RainEntriesMm between two timestamps
	GetRainMMFrom(from time.Time, to time.Time) (*RainEntriesMm, error)

	// GetLastRainTime shows the date of the last rain
	GetLastRainTime() (time.Time, error)

	// GetTempDataCSince gets a TempEntriesC from a time in the past to the present
	GetTempDataCSince(since time.Time) (*TempEntriesC, error)

	// GetTempDataCFrom gets a TempEntriesC between two timestamps
	GetTempDataCFrom(from time.Time, to time.Time) (*TempEntriesC, error)

	// GetLastTempC shows the most recent temperature
	GetLastTempC() (int, error)

	// IsGatewayUp tells whether the gateway has published a status message in a certain time
	IsGatewayUp(since time.Duration) (bool, error)

	// IsSensorUp tells whether the sensor has published a status message in a certain time
	IsSensorUp(since time.Duration) (bool, error)

	// GetEventMessagesSince gets an EventEntries from a time in the past to present. Specify tag or -1 for all tags
	GetEventMessagesSince(tag int, since time.Time) (*EventEntries, error)

	// GetEventMessagesFrom gets an EventEntries between two timestamps. Specify tag or -1 for all tags
	GetEventMessagesFrom(tag int, from time.Time, to time.Time) (*EventEntries, error)

	// Close closes the connection with the database. Necessary for pooled connections
	Close()
}

// RainEntriesMm is a simple array of RainEntryMm values
type RainEntriesMm []RainEntryMm

// RainEntryMm is a single timestamp/mm of rain entry
type RainEntryMm struct {
	Timestamp   time.Time // timestamp on the gateway that the event was recorded
	Millimeters float64   // amount of rain in millimeters
}

// TempEntriesC is an ordered slice of TempEntryC values
type TempEntriesC []TempEntryC

// TempEntryC is a single temperature/timestamp entry
type TempEntryC struct {
	Timestamp time.Time // timestamp on the gateway that the measurement was recorded
	TempC     int       // temperature value in Celsius
}

// EventEntries is a slice of EventEntry structs
type EventEntries []EventEntry

// EventEntry is a single sensor event
type EventEntry struct {
	Timestamp time.Time // timestamp on the gateway that the event was recorded
	Tag       int       // type of event
	Value     int       // value of the event, basically 1 for all events
	Longname  string    // human-comprehensible name, matches 1-to-1 with Tag
}
