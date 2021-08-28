package localdb

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite" // driver for localdb
)

/* Wrap queries in methods so we don't expose the actual databse to the rest of the application */

const (
	foreignKey = `PRAGMA foreign_keys = ON;`
	sqlite     = "sqlite"
)

// connection gets a DB and Conn struct for a sqlite file
type connection struct {
	database *sql.DB
	conn     *sql.Conn
}

type Sqlite struct {
	file     *os.File        // pointer to actual file
	fullPath string          // full POSIX path of sqlite file
	driver   string          // change the type of postgresql connection
	ctx      context.Context // background context
}

// NewSqlite makes a new connector struct for localdb
func NewSqlite(fullPath string, clobber bool) (*Sqlite, error) {
	logrus.Debug("making new Sqlite")
	if clobber {
		_ = os.Remove(fullPath)
	}

	// connect to the file and open it
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}

	// make a Sqlite object and make the schema if necessary
	db := Sqlite{
		file:     file,
		fullPath: fullPath,
		driver:   sqlite,
		ctx:      context.Background(),
	}
	if clobber {
		_, err = db.MakeSchema()
		if err != nil {
			return nil, err
		}
	}
	return &db, nil
}

func (db *Sqlite) newConnection() (*connection, error) {
	// get variables ready
	var (
		database *sql.DB
		conn     *sql.Conn
		err      error
	)

	switch db.driver {
	case sqlite:
		database, err = sql.Open("sqlite", db.fullPath)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		// make a Conn
		conn, err = database.Conn(db.ctx)
		if err != nil {
			logrus.Error("unable to get a connection struct")
			return nil, err
		}
	default:
		panic("unsupported")
	}
	return &connection{database, conn}, nil
}

func (c *connection) disconnect() {
	if err := c.conn.Close(); err != nil {
		logrus.Error(err)
	}
	if err := c.database.Close(); err != nil {
		logrus.Error(err)
	}
}

func (db *Sqlite) MakeSchema() (sql.Result, error) {
	return db.EnterData(localDbSchema)
}

func (db *Sqlite) EnterData(cmd string) (sql.Result, error) {
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

func (db *Sqlite) AddRecord(tag, value int) (sql.Result, error) {
	timestamp := time.Now().Format(time.RFC3339)
	cmd := fmt.Sprintf("INSERT INTO log (tag, value, timestamp) VALUES (%d, %d, \"%s\");", tag, value, timestamp)
	return db.EnterData(cmd)
}

func (db *Sqlite) Tally(tag int) int {
	query := fmt.Sprintf("SELECT COUNT(*) FROM log WHERE tag = %d;", tag)
	return db.GetSingleInt(query)
}

func (db *Sqlite) GetLastRecord(tag int) int {
	cmd := fmt.Sprintf(`SELECT value FROM log WHERE tag = %d ORDER BY id DESC LIMIT 1;`, tag)
	return db.GetSingleInt(cmd)
}

func (db *Sqlite) GetSingleInt(query string) int {
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

// ForeignKeysAreImplemented tests function to ensure foreign key implementation
func (db *Sqlite) ForeignKeysAreImplemented() bool {
	illegal := `INSERT INTO log (tag, value, timestamp) VALUES (99999,1,"timestamp");`
	res, err := db.EnterData(illegal)
	return res == nil && err != nil
}
