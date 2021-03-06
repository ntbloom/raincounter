package database

import (
	"context"
	"database/sql"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite" // Driver for localdb
)

const (
	foreignKey   = `PRAGMA foreign_keys = ON;`
	sqliteDriver = "sqlite"
)

// Connection gets a DB and Conn struct for a sqlite File
type Connection struct {
	Database *sql.DB   // database struct
	Conn     *sql.Conn // connection struct
}

// Sqlite handles connections to sqlite database
type Sqlite struct {
	File     *os.File // name of the .db file
	FullPath string   // full POSIX path
	Driver   string   // sqlite driver
}

// NewSqlite makes a new connector struct for any sqlite database
func NewSqlite(fullPath string, clobber bool, schema string) (*Sqlite, error) {
	if clobber {
		_ = os.Remove(fullPath)
	}

	// connect to the File and open it
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}

	// make a LocalDB object and make the schema if necessary
	db := Sqlite{
		File:     file,
		FullPath: fullPath,
		Driver:   sqliteDriver,
	}
	if clobber {
		_, err = db.MakeSchema(schema)
		if err != nil {
			return nil, err
		}
	}
	return &db, nil
}

// Connect attaches to the sqlite database
func (db *Sqlite) Connect() (*Connection, error) {
	// get variables ready
	var (
		dbPtr *sql.DB
		conn  *sql.Conn
		err   error
	)

	switch db.Driver {
	case sqliteDriver:
		dbPtr, err = sql.Open(sqliteDriver, db.FullPath)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		// make a Conn
		conn, err = dbPtr.Conn(context.Background())
		if err != nil {
			logrus.Error("unable to get a connection struct")
			return nil, err
		}
	default:
		panic("unsupported")
	}
	return &Connection{dbPtr, conn}, nil
}

// Disconnect closes the connection to the database
func (c *Connection) Disconnect() {
	if err := c.Conn.Close(); err != nil {
		logrus.Error(err)
	}
	if err := c.Database.Close(); err != nil {
		logrus.Error(err)
	}
}

// MakeSchema creates all of the database tables, etc.
func (db *Sqlite) MakeSchema(schema string) (sql.Result, error) {
	return db.EnterData(schema)
}

// EnterData enters a generic sql INSERT statement
func (db *Sqlite) EnterData(cmd string) (sql.Result, error) {
	var c *Connection
	var err error

	// enforce foreign keys
	safeCmd := strings.Join([]string{foreignKey, cmd}, " ")
	if c, err = db.Connect(); err != nil {
		return nil, err
	}
	defer c.Disconnect()

	return c.Conn.ExecContext(context.Background(), safeCmd)
}
