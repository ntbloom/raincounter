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
	entry   webdb.DBEntry
	query   webdb.DBQuery
	rainAmt float32
}

func TestWebDB(t *testing.T) {
	test := new(WebDBTest)
	suite.Run(t, test)
}

func (suite *WebDBTest) SetupSuite() {
	config.Configure()

	var entry webdb.DBEntry
	var query webdb.DBQuery
	db, err := webdb.NewPGConnector("raincounter")
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

func (suite *WebDBTest) TestPosgresqlConnection() {
	res, err := suite.query.RunCmd("SELECT 2+2;")
	if err != nil {
		suite.Fail("problem running query", err)
	}
	assert.Equal(suite.T(), 4, suite.query.Unwrap(res), "2+2 != 4")

}

//func (suite *WebDBTest) TestEnterRainEvent() {
//	start := time.Now()
//	time.Sleep(time.Second)
//	timestamp := time.Now().String()
//	qty := 10
//	for i := 0; i < qty; i++ {
//		_, err := suite.entry.AddRainEvent(suite.rainAmt, timestamp)
//		if err != nil {
//			panic(err)
//		}
//	}
//	assert.InDelta(suite.T(), suite.query.TallyRainSince(start), float32(qty)*suite.rainAmt, 0.0001)
//}
