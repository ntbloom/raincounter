package database

// Prep a postgresql.  This is essentially a test fixture but also to be called at the start of a new deployment
import (
	"database/sql"

	"github.com/ntbloom/raincounter/pkg/gateway/tlv"

	"github.com/sirupsen/logrus"
)

// DBWrapper abstracts the underlying SQL engine to let us use the same
// code for localdb, postgresql, or other postgresql
type DBWrapper interface {
	// MakeSchema initializes a schema
	MakeSchema() (sql.Result, error)

	// EnterData enters data into the postgresql without returning any rows
	EnterData(cmd string) (sql.Result, error)

	// AddRecord makes a single integer entry into the postgresql for a given tag
	AddRecord(tag, value int) (sql.Result, error)

	// Tally runs sql command to count postgresql entries for a given topic
	Tally(tag int) int

	// GetLastRecord gets the last record for a given tag
	GetLastRecord(tag int) int

	// GetSingleInt returns the first result of any SQL query that gives at least one integer result
	// simple function for confirming correct value was entered for, say, temperature
	GetSingleInt(query string) int
}

// MakeRainEntry AddRecord a rain event
func MakeRainEntry(db DBWrapper) {
	_, err := db.AddRecord(tlv.Rain, tlv.RainValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeSoftResetEntry AddRecord a soft reset event
func MakeSoftResetEntry(db DBWrapper) {
	_, err := db.AddRecord(tlv.SoftReset, tlv.SoftResetValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeHardResetEntry AddRecord a hard reset event
func MakeHardResetEntry(db DBWrapper) {
	_, err := db.AddRecord(tlv.HardReset, tlv.HardResetValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakePauseEntry AddRecord a pause event
func MakePauseEntry(db DBWrapper) {
	_, err := db.AddRecord(tlv.Pause, tlv.Unpause)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeUnpauseEntry AddRecord an unpause event
func MakeUnpauseEntry(db DBWrapper) {
	_, err := db.AddRecord(tlv.Unpause, tlv.UnpauseValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeTemperatureEntry AddRecord a temperature measurement
func MakeTemperatureEntry(db DBWrapper, tempC int) {
	_, err := db.AddRecord(tlv.Temperature, tempC)
	if err != nil {
		logrus.Error(err)
	}
}

/* GETTERS, MOSTLY FOR TESTING */

func GetRainEntries(db DBWrapper) int {
	return db.Tally(tlv.Rain)
}

func GetSoftResetEntries(db DBWrapper) int {
	return db.Tally(tlv.SoftReset)
}

func GetHardResetEntries(db DBWrapper) int {
	return db.Tally(tlv.HardReset)
}

func GetPauseEntries(db DBWrapper) int {
	return db.Tally(tlv.Pause)
}

func GetUnpauseEntries(db DBWrapper) int {
	return db.Tally(tlv.Unpause)
}

// GetLastTemperatureEntry returns last temp reading, sorted by primary key
func GetLastTemperatureEntry(db DBWrapper) int {
	return db.GetLastRecord(tlv.Temperature)
}
