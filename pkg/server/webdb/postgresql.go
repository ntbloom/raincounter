package webdb

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/ntbloom/raincounter/pkg/common/exitcodes"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/spf13/viper"

	"github.com/ntbloom/raincounter/pkg/gateway/tlv"
	"github.com/sirupsen/logrus"
)

// PGConnector implements the database wrapper used on the cloud server.
// We still use the same tag structure, but put the data into different
// tables because we don't really care about events, just the rain and
// temperature data.

type PGConnector struct {
	ctx  context.Context
	url  string
	pool *pgxpool.Pool
}

func NewPGConnector() *PGConnector {
	ctx := context.Background()
	dbName := viper.GetString(configkey.DatabaseRemoteName)
	password := viper.GetString(configkey.DatabasePostgresqlPassword)
	url := fmt.Sprintf("postgresql://postgres:%s@127.0.0.1:5432/%s", password, dbName)
	logrus.Debugf("connecting to postgres: %s", url)

	duration := time.Millisecond * 200
	waits := int((time.Second * 10) / duration)
	var pgpool *pgxpool.Pool
	var err error
	for i := 0; i < waits; i++ {
		pgpool, err = pgxpool.Connect(ctx, url)
		if err == nil {
			break
		}
		time.Sleep(duration)
	}
	if err != nil {
		logrus.Fatal(err)
		os.Exit(exitcodes.PostgresqlConnnectionError)
	}
	return &PGConnector{ctx, url, pgpool}
}

func (pg *PGConnector) Close() {
	logrus.Info("closing connection pool to postgresql")
	pg.pool.Close()
}

func (pg *PGConnector) Insert(cmd string) error {
	panic("implement me!")
}

func (pg *PGConnector) Select(cmd string) (interface{}, error) {
	return pg.pool.Query(pg.ctx, cmd)
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

func (pg *PGConnector) tallyFloat(table string) float32 {
	panic("implement me!")
}
