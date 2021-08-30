package webdb_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/ntbloom/raincounter/pkg/server/webdb"
	"github.com/spf13/viper"

	"github.com/ntbloom/raincounter/pkg/config"
	"github.com/stretchr/testify/suite"
)

type WebDBTest struct {
	suite.Suite
	entry   webdb.DBEntry
	query   webdb.DBQuery
	rainAmt float32
}

func TestReceiver(t *testing.T) {
	test := new(WebDBTest)
	suite.Run(t, test)
}

func (suite *WebDBTest) SetupSuite() {
	config.Configure()
	sqliteFile := viper.GetString(configkey.DatabaseRemoteFile)

	// use the same object as entry and query
	// in production you would have two separate objects
	var entry webdb.DBEntry
	var query webdb.DBQuery
	db, err := webdb.NewWebSqlite(sqliteFile, true)
	if err != nil {
		panic(err)
	}
	entry = db
	query = db
	suite.entry = entry
	suite.query = query
	suite.rainAmt = float32(viper.GetFloat64(configkey.SensorRainMm))
}
func (suite *WebDBTest) TearDownSuite() {}
func (suite *WebDBTest) SetupTest()     {}
func (suite *WebDBTest) TearDownTest()  {}

func (suite *WebDBTest) TestEnterRainEvent() {
	start := time.Now()
	time.Sleep(time.Second)
	timestamp := time.Now().String()
	qty := 10
	for i := 0; i < qty; i++ {
		_, err := suite.entry.AddRainEvent(suite.rainAmt, timestamp)
		if err != nil {
			panic(err)
		}
	}
	assert.InDelta(suite.T(), suite.query.TallyRainSince(start), float32(qty)*suite.rainAmt, 0.0001)
}
