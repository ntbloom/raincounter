package database

// Prep a database.  This is essentially a test fixture but also to be called at the start of a new deployment
import (
	"database/sql"

	tlv2 "github.com/ntbloom/raincounter/pkg/gateway/tlv"

	"github.com/sirupsen/logrus"
)

// DBWrapper abstracts the underlying SQL engine to let us use the same
// code for sqlite, postgresql, or other database
type DBWrapper interface {
	// MakeSchema initializes a schema
	MakeSchema() (sql.Result, error)

	// EnterData enters data into the database without returning any rows
	EnterData(cmd string) (sql.Result, error)

	// AddRecord makes a single integer entry into the database for a given tag
	AddRecord(tag, value int) (sql.Result, error)

	// Tally runs sql command to count database entries for a given topic
	Tally(tag int) int

	// GetLastRecord gets the last record for a given tag
	GetLastRecord(tag int) int

	// GetSingleInt returns the first result of any SQL query that gives at least one integer result
	// simple function for confirming correct value was entered for, say, temperature
	GetSingleInt(query string) int
}

// MakeRainEntry AddRecord a rain event
func MakeRainEntry(db DBWrapper) {
	_, err := db.AddRecord(tlv2.Rain, tlv2.RainValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeSoftResetEntry AddRecord a soft reset event
func MakeSoftResetEntry(db DBWrapper) {
	_, err := db.AddRecord(tlv2.SoftReset, tlv2.SoftResetValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeHardResetEntry AddRecord a hard reset event
func MakeHardResetEntry(db DBWrapper) {
	_, err := db.AddRecord(tlv2.HardReset, tlv2.HardResetValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakePauseEntry AddRecord a pause event
func MakePauseEntry(db DBWrapper) {
	_, err := db.AddRecord(tlv2.Pause, tlv2.Unpause)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeUnpauseEntry AddRecord an unpause event
func MakeUnpauseEntry(db DBWrapper) {
	_, err := db.AddRecord(tlv2.Unpause, tlv2.UnpauseValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeTemperatureEntry AddRecord a temperature measurement
func MakeTemperatureEntry(db DBWrapper, tempC int) {
	_, err := db.AddRecord(tlv2.Temperature, tempC)
	if err != nil {
		logrus.Error(err)
	}
}

/* GETTERS, MOSTLY FOR TESTING */

func GetRainEntries(db DBWrapper) int {
	return db.Tally(tlv2.Rain)
}

func GetSoftResetEntries(db DBWrapper) int {
	return db.Tally(tlv2.SoftReset)
}

func GetHardResetEntries(db DBWrapper) int {
	return db.Tally(tlv2.HardReset)
}

func GetPauseEntries(db DBWrapper) int {
	return db.Tally(tlv2.Pause)
}

func GetUnpauseEntries(db DBWrapper) int {
	return db.Tally(tlv2.Unpause)
}

// GetLastTemperatureEntry returns last temp reading, sorted by primary key
func GetLastTemperatureEntry(db DBWrapper) int {
	return db.GetLastRecord(tlv2.Temperature)
}
