package database

import (
	"database/sql"

	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite" // driver for sqlite
)

/* Wrap queries in methods so we don't expose the actual databse to the rest of the application */

// connection gets a DB and Conn struct for a sqlite file
type connection struct {
	database *sql.DB
	conn     *sql.Conn
}

func (db *DBConnector) newConnection() (*connection, error) {
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
			logrus.Error("unable to open database")
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
