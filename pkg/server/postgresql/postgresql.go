// postgresql runs the postgresql code
package postgresql

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type PgConnector struct {
	dbName string // name of the postgresql
	url    string // how to connect

}

// NewDatabase is an empty function but will eventually do something
func NewDatabase(dbName, url string) *PgConnector {
	logrus.Debugf("this is the sql schema: %s", sqlSchema)
	return &PgConnector{
		dbName: dbName,
		url:    url,
	}
}

// tests the connection, can delete later
func (p *PgConnector) MakeContact() error {
	conn, err := pgxpool.Connect(context.Background(), p.url)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer conn.Close()
	return nil
}
