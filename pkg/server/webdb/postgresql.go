package webdb

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ntbloom/raincounter/pkg/gateway/tlv"

	"github.com/jackc/pgx/v4"

	"github.com/ntbloom/raincounter/pkg/common/exitcodes"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

// PGConnector implements the database wrapper used on the cloud server.
// We still use the same tag structure, but put the data into different
// tables because we don't really care about events, just the rain and
// temperature data.

const (
	intErrVal   = -999
	floatErrVal = -999.0
)

var errTime time.Time = time.Unix(0, 0)

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
	if err != nil {
		logrus.Error(err)
	}
	defer res.Close()
	return err
}

func (pg *PGConnector) AddTagValue(tag int, value int, gwTimestamp time.Time) error {
	switch tag {
	// don't use these methods
	case tlv.Rain:
		return fmt.Errorf("rain events not supported in AddTagValue")
	case tlv.Temperature:
		return fmt.Errorf("temperature events not supported in AddTagValue")
	default:
		sql := fmt.Sprintf(
			`INSERT INTO event_log (gw_timestamp, server_timestamp, tag, value) VALUES ('%s','%s',%d,%d);`,
			gwTimestamp.Format(configkey.TimestampFormat), time.Now().Format(configkey.TimestampFormat), tag, value,
		)
		return pg.Insert(sql)
	}
}

func (pg *PGConnector) AddStatusUpdate(asset int, gwTimestamp time.Time) error {
	sql := fmt.Sprintf(`INSERT INTO status_log (gw_timestamp, server_timestamp, asset) VALUES ('%s','%s',%d);`,
		gwTimestamp.Format(configkey.TimestampFormat), time.Now().Format(configkey.TimestampFormat), asset)
	return pg.Insert(sql)
}

func (pg *PGConnector) AddTempCValue(tempC int, gwTimestamp time.Time) error {
	sql := fmt.Sprintf(
		`INSERT INTO temperature (gw_timestamp, server_timestamp, value) VALUES ('%s','%s',%d);`,
		gwTimestamp.Format(configkey.TimestampFormat), time.Now().Format(configkey.TimestampFormat), tempC)
	return pg.Insert(sql)
}

func (pg *PGConnector) AddRainMMEvent(amount float32, gwTimestamp time.Time) error {
	sql := fmt.Sprintf(
		`INSERT INTO rain (gw_timestamp, server_timestamp, amount ) VALUES ('%s','%s','%f');`,
		gwTimestamp.Format(configkey.TimestampFormat), time.Now().Format(configkey.TimestampFormat), amount)
	return pg.Insert(sql)
}

/* QUERYING RAIN */

func (pg *PGConnector) Select(cmd string) (interface{}, error) {
	return pg.genericQuery(cmd)
}

func (pg *PGConnector) TotalRainMMSince(since time.Time) float32 {
	return pg.TotalRainMMFrom(since, time.Now())
}

func (pg *PGConnector) TotalRainMMFrom(from, to time.Time) float32 {
	sql := fmt.Sprintf(`SELECT sum(amount) FROM rain WHERE gw_timestamp BETWEEN '%s' and '%s';`,
		from.Format(configkey.TimestampFormat), to.Format(configkey.TimestampFormat))
	row, err := pg.genericQuery(sql)
	if err != nil {
		logrus.Error(err)
		return floatErrVal
	}
	defer row.Close()
	row.Next()
	var total float32
	err = row.Scan(&total)
	if err != nil {
		logrus.Error(err)
		return floatErrVal
	}
	return total
}

func (pg *PGConnector) GetRainMMSince(since time.Time) *RainEntriesMm {
	return pg.GetRainMMFrom(since, time.Now())
}

func (pg *PGConnector) GetRainMMFrom(from, to time.Time) *RainEntriesMm {
	sql := fmt.Sprintf(`
		SELECT gw_timestamp, amount 
		FROM rain 
		WHERE gw_timestamp BETWEEN '%s' and '%s'
		ORDER BY gw_timestamp
		;
	`, from.Format(configkey.TimestampFormat), to.Format(configkey.TimestampFormat))
	rows, err := pg.genericQuery(sql)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	defer rows.Close()
	var rain RainEntriesMm
	for rows.Next() {
		var amt float32
		var stamp time.Time
		err = rows.Scan(&stamp, &amt)
		if err != nil {
			logrus.Error(err)
			return nil
		}
		rain = append(rain, RainEntryMm{
			Timestamp:   stamp,
			Millimeters: amt,
		})
	}
	return &rain
}

func (pg *PGConnector) GetLastRainTime() time.Time {
	sql := `SELECT amount FROM rain ORDER BY gw_timestamp DESC LIMIT 1;`
	row, err := pg.genericQuery(sql)
	if err != nil {
		return errTime
	}
	defer row.Close()
	var stamp time.Time
	row.Next()
	err = row.Scan(&stamp)
	if err != nil {
		logrus.Errorf("failure to scan row for last rain timestamp: %s", err)
		return errTime
	}
	return stamp
}

/* QUERYING TEMPERATURE */

func (pg *PGConnector) GetTempDataCSince(since time.Time) *TempEntriesC {
	return pg.GetTempDataCFrom(since, time.Now())
}

func (pg *PGConnector) GetTempDataCFrom(from time.Time, to time.Time) *TempEntriesC {
	sql := fmt.Sprintf(`
		SELECT gw_timestamp, value
		FROM temperature
		WHERE gw_timestamp BETWEEN '%s' and '%s'
		ORDER BY gw_timestamp
		;
	`, from.Format(configkey.TimestampFormat), to.Format(configkey.TimestampFormat))
	rows, err := pg.genericQuery(sql)
	if err != nil {
		logrus.Errorf("bad query: `%s`", sql)
		return nil
	}
	defer rows.Close()

	var temps TempEntriesC
	for rows.Next() {
		var timestamp time.Time
		var tempC int
		err := rows.Scan(&timestamp, &tempC)
		if err != nil {
			logrus.Errorf("cannot retrieve timestamp/tempC row: %s", err)
			return nil
		}
		temps = append(temps, TempEntryC{
			timestamp,
			tempC,
		})
	}
	return &temps
}

func (pg *PGConnector) GetLastTempC() int {
	sql := `SELECT value FROM temperature ORDER BY gw_timestamp DESC LIMIT 1;`
	row, err := pg.genericQuery(sql)
	if err != nil {
		logrus.Error(err)
		return intErrVal
	}
	defer row.Close()
	var tempC int
	row.Next()
	err = row.Scan(&tempC)
	if err != nil {
		logrus.Errorf("failed to scan row for tempC: %s", err)
		return intErrVal
	}
	return tempC
}

/* RANDOM HELPER FUNCTIONS */

// executes arbitrary sql. we need to close the connection after each value, either for
func (pg *PGConnector) genericQuery(cmd string) (pgx.Rows, error) {
	logrus.Debugf("pgsql: %s", cmd)
	return pg.pool.Query(pg.ctx, cmd)
}
