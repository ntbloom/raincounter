package webdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/spf13/viper"

	"github.com/jackc/pgx/v4"
	"github.com/ntbloom/raincounter/pkg/gateway/tlv"
	"github.com/sirupsen/logrus"
)

// PGConnector implements the database wrapper used on the cloud server.
// We still use the same tag structure, but put the data into different
// tables because we don't really care about events, just the rain and
// temperature data.

type PGConnector struct {
	ctx context.Context
	url string
}

func NewPGConnector() *PGConnector {
	ctx := context.Background()
	dbName := viper.GetString(configkey.DatabaseRemoteName)
	password := viper.GetString(configkey.DatabasePostgresqlPassword)
	url := fmt.Sprintf("postgresql://postgres:%s@127.0.0.1:5432/%s", password, dbName)
	logrus.Error(url)

	return &PGConnector{ctx, url}
}

func (pg *PGConnector) RunCmd(cmd string) (sql.Result, error) {
	conn, err := pg.connect()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		err := conn.Close(pg.ctx)
		if err != nil {
			logrus.Warningf("connection not closed properly: %s", err)
		}
	}()
	return nil, nil
}

func (pg *PGConnector) Unwrap(sql.Result) interface{} {
	panic("implement me!")
}

func (pg *PGConnector) AddTagValue(tag, value int) (sql.Result, error) {
	switch tag {
	// don't use these methods
	case tlv.Rain:
		panic("rain events not supported in this method")
	case tlv.Temperature:
		panic("temperature events not supported in this method")

	case tlv.SoftReset:
		logrus.Debug("adding soft reset to web database")
	case tlv.HardReset:
		logrus.Debug("adding hard reset to web database")
	case tlv.Pause:
		logrus.Debug("adding pause to web database")
	case tlv.Unpause:
		logrus.Debug("adding unpause to web database")
	default:
		panic("unsupported tag")
	}
	return nil, nil
}

func (pg *PGConnector) AddRainEvent(value float32, gwTimestamp string) (sql.Result, error) {
	panic("implement me!")
}

func (pg *PGConnector) TallyRainSince(since time.Time) float32 {
	panic("implement me!")
}

func (pg *PGConnector) TallyRainFrom(start, finish time.Time) float32 {
	panic("implement me!")
}

func (pg *PGConnector) GetLastRainTime() time.Time {
	panic("implement me!")
}

func (pg *PGConnector) connect() (*pgx.Conn, error) {
	return pgx.Connect(pg.ctx, pg.url)
}

func (pg *PGConnector) tallyFloat(table string) float32 {
	panic("implement me!")
}
