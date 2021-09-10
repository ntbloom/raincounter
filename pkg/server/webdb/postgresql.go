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

const ()

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

func (pg *PGConnector) AddRainMMEvent(amount float64, gwTimestamp time.Time) error {
	sql := fmt.Sprintf(
		`INSERT INTO rain (gw_timestamp, server_timestamp, amount ) VALUES ('%s','%s','%f');`,
		gwTimestamp.Format(configkey.TimestampFormat), time.Now().Format(configkey.TimestampFormat), amount)
	return pg.Insert(sql)
}

/* QUERYING RAIN */

func (pg *PGConnector) Select(cmd string) (interface{}, error) {
	return pg.genericQuery(cmd)
}

func (pg *PGConnector) TotalRainMMSince(since time.Time) (float64, error) {
	return pg.TotalRainMMFrom(since, time.Now())
}

func (pg *PGConnector) TotalRainMMFrom(from, to time.Time) (float64, error) {
	sql := fmt.Sprintf(`SELECT sum(amount) FROM rain WHERE gw_timestamp BETWEEN '%s' and '%s';`,
		from.Format(configkey.TimestampFormat), to.Format(configkey.TimestampFormat))
	row, err := pg.genericQuery(sql)
	if err != nil {
		logrus.Error(err)
		return configkey.FloatErrVal, err
	}
	defer row.Close()
	row.Next()
	var total float64
	err = row.Scan(&total)
	if err != nil {
		// means there is no value
		return 0.0, nil
	}
	return total, nil
}

func (pg *PGConnector) GetRainMMSince(since time.Time) (*RainEntriesMm, error) {
	return pg.GetRainMMFrom(since, time.Now())
}

func (pg *PGConnector) GetRainMMFrom(from, to time.Time) (*RainEntriesMm, error) {
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
		return nil, err
	}
	defer rows.Close()
	var rain RainEntriesMm
	for rows.Next() {
		var amt float64
		var stamp time.Time
		err = rows.Scan(&stamp, &amt)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		rain = append(rain, RainEntryMm{
			Timestamp:   stamp,
			Millimeters: amt,
		})
	}
	return &rain, nil
}

func (pg *PGConnector) GetLastRainTime() (time.Time, error) {
	sql := `SELECT gw_timestamp FROM rain ORDER BY gw_timestamp DESC LIMIT 1;`
	row, err := pg.genericQuery(sql)
	if err != nil {
		return errTime, err
	}
	defer row.Close()
	var stamp time.Time
	row.Next()
	err = row.Scan(&stamp)
	if err != nil {
		logrus.Errorf("failure to scan row for last rain timestamp: %s", err)
		return errTime, err
	}
	return stamp, nil
}

/* QUERYING TEMPERATURE */

func (pg *PGConnector) GetTempDataCSince(since time.Time) (*TempEntriesC, error) {
	return pg.GetTempDataCFrom(since, time.Now())
}

func (pg *PGConnector) GetTempDataCFrom(from time.Time, to time.Time) (*TempEntriesC, error) {
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
		return nil, err
	}
	defer rows.Close()

	var temps TempEntriesC
	for rows.Next() {
		var timestamp time.Time
		var tempC int
		err := rows.Scan(&timestamp, &tempC)
		if err != nil {
			logrus.Errorf("cannot retrieve timestamp/tempC row: %s", err)
			return nil, err
		}
		temps = append(temps, TempEntryC{
			timestamp,
			tempC,
		})
	}
	return &temps, nil
}

func (pg *PGConnector) GetLastTempC() (int, error) {
	sql := `SELECT value FROM temperature ORDER BY gw_timestamp DESC LIMIT 1;`
	row, err := pg.genericQuery(sql)
	if err != nil {
		logrus.Error(err)
		return configkey.IntErrVal, err
	}
	defer row.Close()
	var tempC int
	row.Next()
	err = row.Scan(&tempC)
	if err != nil {
		logrus.Errorf("failed to scan row for tempC: %s", err)
		return configkey.IntErrVal, err
	}
	return tempC, nil
}

func (pg *PGConnector) IsGatewayUp(since time.Duration) (bool, error) {
	return pg.getLastStatusMessage(since, "gateway")
}

func (pg *PGConnector) IsSensorUp(since time.Duration) (bool, error) {
	return pg.getLastStatusMessage(since, "sensor")
}

func (pg *PGConnector) GetEventMessagesSince(tag int, since time.Time) (*EventEntries, error) {
	return pg.GetEventMessagesFrom(tag, since, time.Now())
}

func (pg *PGConnector) GetEventMessagesFrom(tag int, from, to time.Time) (*EventEntries, error) {
	sqlAll := fmt.Sprintf(`
SELECT mappings.longname, event_log.gw_timestamp, event_log.tag, event_log.value
FROM event_log
LEFT JOIN mappings on event_log.tag = mappings.id
WHERE gw_timestamp BETWEEN '%s' and '%s'
ORDER BY gw_timestamp DESC
;`, from.Format(configkey.TimestampFormat), to.Format(configkey.TimestampFormat))
	sqlTag := fmt.Sprintf(`
SELECT mappings.longname, event_log.gw_timestamp, event_log.tag, event_log.value
FROM event_log
LEFT JOIN mappings on event_log.tag = mappings.id
WHERE event_log.tag = %d
AND gw_timestamp BETWEEN '%s' and '%s'
ORDER BY gw_timestamp DESC
`, tag, from.Format(configkey.TimestampFormat), to.Format(configkey.TimestampFormat))

	var query string
	if tag >= 2 && tag <= 5 {
		query = sqlTag
	} else if tag == -1 {
		query = sqlAll
	} else {
		return nil, fmt.Errorf("illegal tag %d", tag)
	}
	rows, err := pg.genericQuery(query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()

	var entries EventEntries
	for rows.Next() {
		var longname string
		var timestamp time.Time
		var tag int
		var value int
		err = rows.Scan(&longname, &timestamp, &tag, &value)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		entry := EventEntry{
			Timestamp: timestamp,
			Tag:       tag,
			Value:     value,
			Longname:  longname,
		}
		entries = append(entries, entry)
	}
	return &entries, nil

}

/* RANDOM HELPER FUNCTIONS */

// executes arbitrary sql. we need to close the connection after each value, either for
func (pg *PGConnector) genericQuery(cmd string) (pgx.Rows, error) {
	logrus.Debugf("pgsql: %s", cmd)
	return pg.pool.Query(pg.ctx, cmd)
}

func (pg *PGConnector) getLastStatusMessage(since time.Duration, asset string) (bool, error) {
	sql := fmt.Sprintf(`
SELECT gw_timestamp 
FROM status_log 
LEFT JOIN status_codes on status_log.asset = status_codes.id 
WHERE status_codes.asset = '%s'
ORDER BY gw_timestamp DESC
LIMIT 1
;`, asset)
	row, err := pg.genericQuery(sql)
	if err != nil {
		logrus.Error(err)
		return false, err
	}
	defer row.Close()
	row.Next()

	var timestamp time.Time
	err = row.Scan(&timestamp)
	if err != nil {
		logrus.Error(err)
		return false, err
	}
	now := time.Now()
	diff := now.Sub(timestamp)
	return diff < since, nil
}
