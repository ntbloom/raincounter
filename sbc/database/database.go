package database

// Prep a database.  This is essentially a test fixture but also to be called at the start of a new deployment
import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ntbloom/rainbase/pkg/tlv"

	"github.com/sirupsen/logrus"
)

const foreignKey = `PRAGMA foreign_keys = ON;`
const sqlite = "sqlite"

type DBConnector struct {
	file     *os.File        // pointer to actual file
	fullPath string          // full POSIX path of sqlite file
	driver   string          // change the type of database connection
	ctx      context.Context // background context
}

// NewSqliteDBConnector makes a new connector struct for sqlite
func NewSqliteDBConnector(fullPath string, clobber bool) (*DBConnector, error) {
	logrus.Debug("making new DBConnector")
	if clobber {
		_ = os.Remove(fullPath)
	}

	// connect to the file and open it
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}

	// make a DBConnector object and make the schema if necessary
	db := DBConnector{
		file:     file,
		fullPath: fullPath,
		driver:   sqlite,
		ctx:      context.Background(),
	}
	if clobber {
		_, err = db.makeSchema()
		if err != nil {
			return nil, err
		}
	}
	return &db, nil
}

/* SQL LOG ENTRIES */
// MakeRainEntry addRecord a rain event
func (db *DBConnector) MakeRainEntry() {
	_, err := db.addRecord(tlv.Rain, tlv.RainValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeSoftResetEntry addRecord a soft reset event
func (db *DBConnector) MakeSoftResetEntry() {
	_, err := db.addRecord(tlv.SoftReset, tlv.SoftResetValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeHardResetEntry addRecord a hard reset event
func (db *DBConnector) MakeHardResetEntry() {
	_, err := db.addRecord(tlv.HardReset, tlv.HardResetValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakePauseEntry addRecord a pause event
func (db *DBConnector) MakePauseEntry() {
	_, err := db.addRecord(tlv.Pause, tlv.Unpause)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeUnpauseEntry addRecord an unpause event
func (db *DBConnector) MakeUnpauseEntry() {
	_, err := db.addRecord(tlv.Unpause, tlv.UnpauseValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeTemperatureEntry addRecord a temperature measurement
func (db *DBConnector) MakeTemperatureEntry(tempC int) {
	_, err := db.addRecord(tlv.Temperature, tempC)
	if err != nil {
		logrus.Error(err)
	}
}

/* GETTERS, MOSTLY FOR TESTING */
func (db *DBConnector) GetRainEntries() int {
	return db.tally(tlv.Rain)
}

func (db *DBConnector) GetSoftResetEntries() int {
	return db.tally(tlv.SoftReset)
}

func (db *DBConnector) GetHardResetEntries() int {
	return db.tally(tlv.HardReset)
}

func (db *DBConnector) GetPauseEntries() int {
	return db.tally(tlv.Pause)
}

func (db *DBConnector) GetUnpauseEntries() int {
	return db.tally(tlv.Unpause)
}

// GetLastTemperatureEntry returns last temp reading, sorted by primary key
func (db *DBConnector) GetLastTemperatureEntry() int {
	return db.getLastRecord(tlv.Temperature)
}

/* SELECTED METHODS EXPORTED FOR TEST/VERIFICATION */

// ForeignKeysAreImplemented, test function to ensure foreign key implementation
func (db *DBConnector) ForeignKeysAreImplemented() bool {
	illegal := `INSERT INTO log (tag, value, timestamp) VALUES (99999,1,"timestamp");`
	res, err := db.enterData(illegal)
	return res == nil && err != nil
}

/* HELPER METHODS */

// makeSchema puts the schema in the sqlite file
func (db *DBConnector) makeSchema() (sql.Result, error) {
	return db.enterData(sqlschema)
}

// enterData enters data into the database without returning any rows
func (db *DBConnector) enterData(cmd string) (sql.Result, error) {
	var c *connection
	var err error

	// enforce foreign keys
	safeCmd := strings.Join([]string{foreignKey, cmd}, " ")
	if c, err = db.newConnection(); err != nil {
		return nil, err
	}
	defer c.disconnect()

	return c.conn.ExecContext(db.ctx, safeCmd)
}

// addRecord makes an entry into the databse
// base command for all logging
func (db *DBConnector) addRecord(tag, value int) (sql.Result, error) {
	timestamp := time.Now().Format(time.RFC3339)
	cmd := fmt.Sprintf("INSERT INTO log (tag, value, timestamp) VALUES (%d, %d, \"%s\");", tag, value, timestamp)
	return db.enterData(cmd)
}

/* QUERYING METHODS, MOSTLY FOR TESTING */

// tally runs sql command to count database entries for a given topic
func (db *DBConnector) tally(tag int) int {
	query := fmt.Sprintf("SELECT COUNT(*) FROM log WHERE tag = %d;", tag)
	return db.getSingleInt(query)
}

// getLastRecord gets the last record for a given tag
func (db *DBConnector) getLastRecord(tag int) int {
	cmd := fmt.Sprintf(`SELECT value FROM log WHERE tag = %d ORDER BY id DESC LIMIT 1;`, tag)
	return db.getSingleInt(cmd)
}

// getSingleInt returns the first result of any SQL query that gives at least one integer result
// simple function for confirming correct value was entered for, say, temperature
func (db *DBConnector) getSingleInt(query string) int {
	var rows *sql.Rows
	var err error

	c, _ := db.newConnection() // don't handle the error, just return -1
	defer c.disconnect()

	if rows, err = c.conn.QueryContext(db.ctx, query); err != nil {
		return -1
	}
	closed := func() {
		if err = rows.Close(); err != nil {
			logrus.Error(err)
		}
	}
	defer closed()
	results := make([]int, 0)
	for rows.Next() {
		var val int
		if err = rows.Scan(&val); err != nil {
			logrus.Error(err)
			return -1
		}
		results = append(results, val)
	}

	return results[0]
}
