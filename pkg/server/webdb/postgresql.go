package webdb

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4"

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
	pool *pgxpool.Pool
}

func NewPGConnector() *PGConnector {
	ctx := context.Background()
	dbName := viper.GetString(configkey.PGDatabaseName)
	password := viper.GetString(configkey.PGPassword)
	url := fmt.Sprintf("postgresql://postgres:%s@127.0.0.1:5432/%s", password, dbName)
	logrus.Debugf("connecting to postgres: %s", url)

	duration := viper.GetDuration(configkey.PGConnectionRetryWait)
	totalWait := int((viper.GetDuration(configkey.PGConnectionTimeout)) / duration)
	var pgpool *pgxpool.Pool
	var err error
	for i := 0; i < totalWait; i++ {
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
	return &PGConnector{ctx, pgpool}
}

func (pg *PGConnector) Close() {
	logrus.Info("closing connection pool to postgresql")
	pg.pool.Close()
}

/* INSERTING DATA */

func (pg *PGConnector) Insert(cmd string) error {
	res, err := pg.genericQuery(cmd)
	res.Close()
	return err
}

func (pg *PGConnector) AddTagValue(tag int, value int, t time.Time) error {
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
	return nil
}

func (pg *PGConnector) AddTempCValue(tempC int, gwTimestamp time.Time) error {
	panic("implement me!")
}

func (pg *PGConnector) AddRainMMEvent(value float32, gwTimestamp time.Time) error {
	panic("implement me!")
}

/* QUERYING DATA */

func (pg *PGConnector) Select(cmd string) (interface{}, error) {
	return pg.genericQuery(cmd)
}

func (pg *PGConnector) TotalRainMMSince(since time.Time) float32 {
	panic("implement me!")
}

func (pg *PGConnector) TotalRainMMFrom(from, to time.Time) float32 {
	panic("implement me!")
}

func (pg *PGConnector) GetRainMMSince(timestamp time.Time) RainMMMap {
	panic("implement me!")
}

func (pg *PGConnector) GetRainMMFrom(from, to time.Time) RainMMMap {
	panic("implement me!")
}

func (pg *PGConnector) GetLastRainTime() time.Time {
	panic("implement me!")
}

/* RANDOM HELPER FUNCTIONS */

// executes arbitrary sql. we need to close the connection after each value, either for
func (pg *PGConnector) genericQuery(cmd string) (pgx.Rows, error) {
	logrus.Debugf("pgsql: %s", cmd)
	return pg.pool.Query(pg.ctx, cmd)
}
