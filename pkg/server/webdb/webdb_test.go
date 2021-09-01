package webdb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jackc/pgx/v4"

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
	db := webdb.NewPGConnector()
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
	res, err := suite.query.Select("SELECT 2+2;")
	if err != nil || res == nil {
		suite.Fail("problem running query", err)
	}
	var sum int
	val := res.(pgx.Rows)
	defer val.Close()
	val.Next()
	err = val.Scan(&sum)
	if err != nil {
		suite.Fail("bad reflection", err)
	}
	assert.Equal(suite.T(), 4, sum, "failed simple SQL math")
}

// func (suite *WebDBTest) TestEnterRainEvent() {
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
