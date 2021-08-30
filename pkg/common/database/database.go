package database

// Prep a postgresql.  This is essentially a test fixture but also to be called at the start of a new deployment
import (
	"database/sql"

	"github.com/ntbloom/raincounter/pkg/gateway/tlv"

	"github.com/sirupsen/logrus"
)

// DBWrapper abstracts the underlying SQL engine/implementations
type DBWrapper interface {
	// MakeSchema initializes a schema
	MakeSchema() (sql.Result, error)

	// EnterData enters data into the database without returning any rows
	EnterData(cmd string) (sql.Result, error)

	// AddIntRecord makes a single integer entry into the database for a given tag
	AddIntRecord(tag, value int) (sql.Result, error)

	// AddFloatRecord makes a single float entry into the database for a given tag
	AddFloatRecord(tag int, value float32) (sql.Result, error)

	// Tally runs sql command to count database entries for a given topic
	Tally(tag int) int

	// GetLastRecord gets the last record for a given tag
	GetLastRecord(tag int) int

	// GetSingleInt returns the first result of any SQL query that gives at least one integer result
	// simple function for confirming correct value was entered for, say, temperature
	GetSingleInt(query string) int
}

// MakeRainTallyEntry AddTag a rain event
func MakeRainTallyEntry(db DBWrapper) {
	_, err := db.AddIntRecord(tlv.Rain, tlv.RainValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeRainValueEntry AddFloatRecord for a rain event
func MakeRainValueEntry(db DBWrapper, value float32) {
	_, err := db.AddFloatRecord(tlv.RainValue, value)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeSoftResetEntry AddTag a soft reset event
func MakeSoftResetEntry(db DBWrapper) {
	_, err := db.AddIntRecord(tlv.SoftReset, tlv.SoftResetValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeHardResetEntry AddTag a hard reset event
func MakeHardResetEntry(db DBWrapper) {
	_, err := db.AddIntRecord(tlv.HardReset, tlv.HardResetValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakePauseEntry AddTag a pause event
func MakePauseEntry(db DBWrapper) {
	_, err := db.AddIntRecord(tlv.Pause, tlv.Unpause)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeUnpauseEntry AddTag an unpause event
func MakeUnpauseEntry(db DBWrapper) {
	_, err := db.AddIntRecord(tlv.Unpause, tlv.UnpauseValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeTemperatureEntry AddTag a temperature measurement
func MakeTemperatureEntry(db DBWrapper, tempC int) {
	_, err := db.AddIntRecord(tlv.Temperature, tempC)
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
