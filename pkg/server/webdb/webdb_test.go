package webdb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/ntbloom/raincounter/pkg/server/webdb"
	"github.com/spf13/viper"

	"github.com/ntbloom/raincounter/pkg/config"
	"github.com/stretchr/testify/suite"
)

type WebDBTest struct {
	suite.Suite
	db      *webdb.WebDB
	rainAmt float32
}

func TestReceiver(t *testing.T) {
	test := new(WebDBTest)
	suite.Run(t, test)
}

func (suite *WebDBTest) SetupSuite() {
	config.Configure()
	sqliteFile := viper.GetString(configkey.DatabaseRemoteFile)
	db, err := webdb.NewWebDB(sqliteFile, true)
	if err != nil {
		panic(err)
	}
	suite.db = db
	suite.rainAmt = float32(viper.GetFloat64(configkey.SensorRainMm))
}
func (suite *WebDBTest) TearDownSuite() {}
func (suite *WebDBTest) SetupTest()     {}
func (suite *WebDBTest) TearDownTest()  {}

func (suite *WebDBTest) TestBasicSchema() {
	qty := 10
	for i := 0; i < qty; i++ {
		_, err := suite.db.AddRainEvent(suite.rainAmt, "timestamp")
		if err != nil {
			panic(err)
		}
	}
	assert.InDelta(suite.T(), suite.db.TallyRain(), float32(qty)*suite.rainAmt, 0.0001)
}
