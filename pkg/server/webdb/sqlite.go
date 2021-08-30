package webdb

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ntbloom/raincounter/pkg/gateway/tlv"
	"github.com/sirupsen/logrus"

	"github.com/ntbloom/raincounter/pkg/common/database"
)

// WebSqlite implements the database wrapper used on the cloud server.
// We still use the same tag structure, but put the data into different
// tables because we don't really care about events, just the rain and
// temperature data.
type WebSqlite struct {
	lite *database.Sqlite
}

func NewWebSqlite(fullPath string, clobber bool) (*WebSqlite, error) {
	lite, err := database.NewSqlite(fullPath, clobber, webDBSchema)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &WebSqlite{lite}, nil
}

func (w *WebSqlite) EnterData(cmd string) (sql.Result, error) {
	return w.lite.EnterData(cmd)
}

func (w *WebSqlite) AddTagValue(tag, value int) (sql.Result, error) {
	switch tag {
	case tlv.Rain:
		panic("rain events not supported in this method")
	case tlv.Temperature:
		logrus.Debug("adding temp record to web database")
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

func (w *WebSqlite) AddRainEvent(value float32, gwTimestamp string) (sql.Result, error) {
	timestamp := time.Now().Format(time.RFC3339)
	cmd := fmt.Sprintf(
		"INSERT INTO rain (gw_timestamp, server_timestamp, amount) VALUES (\"%s\",\"%s\",%f);",
		gwTimestamp,
		timestamp,
		value,
	)
	res, err := w.lite.EnterData(cmd)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return res, nil
}

func (w *WebSqlite) TallyRainSince(since time.Time) float32 {
	panic("implement me!")
	return w.tallyFloat("rain")
}

func (w *WebSqlite) TallyRainFrom(start, finish time.Time) float32 {
	panic("implement me")
}

func (w *WebSqlite) GetLastRain() time.Time {
	panic("implement me!")
}

func (w *WebSqlite) tallyFloat(table string) float32 {
	var rows *sql.Rows
	var err error
	c, _ := w.lite.Connect()
	defer c.Disconnect()

	query := fmt.Sprintf("SELECT SUM(amount) FROM \"%s\"", table)
	if rows, err = c.Conn.QueryContext(w.lite.Ctx, query); err != nil {
		logrus.Error(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	results := make([]float32, 0)
	for rows.Next() {
		var val float32
		if err = rows.Scan(&val); err != nil {
			logrus.Error(err)
			return -1.0
		}
		results = append(results, val)
	}

	return results[0]
}
