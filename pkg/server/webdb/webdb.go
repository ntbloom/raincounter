package webdb

import (
	"database/sql"

	"github.com/ntbloom/raincounter/pkg/gateway/tlv"
	"github.com/sirupsen/logrus"

	"github.com/ntbloom/raincounter/pkg/common/database"
)

// WebDB implements the database wrapper used on the cloud server.
// We still use the same tag structure, but put the data into different
// tables because we don't really care about events, just the rain and
// temperature data.
type WebDB struct {
	lite *database.Sqlite
}

func NewWebDB(fullPath string, clobber bool) (*WebDB, error) {
	lite, err := database.NewSqlite(fullPath, clobber, webDBSchema)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &WebDB{lite}, nil
}

func (w *WebDB) MakeSchema() (sql.Result, error) {
	return w.lite.MakeSchema(webDBSchema)
}

func (w *WebDB) EnterData(cmd string) (sql.Result, error) {
	return w.lite.EnterData(cmd)
}

func (w *WebDB) AddIntRecord(tag, value int) (sql.Result, error) {
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

func (w *WebDB) AddFloatRecord(tag int, value float32) (sql.Result, error) {
	panic("implement me")
}

func (w *WebDB) Tally(tag int) int {
	panic("implement me")
}

func (w *WebDB) TallyRain() float32 {
	return -1.0
}

func (w *WebDB) GetLastRecord(tag int) int {
	panic("implement me")
}

func (w *WebDB) GetSingleInt(query string) int {
	panic("implement me")
}
