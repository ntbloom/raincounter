package webdb_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

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
func (suite *WebDBTest) TearDownSuite() {
	// close the connections
	logrus.Debug("closing the test suite's entry pool...")
	suite.entry.Close()
	logrus.Debug("closing the test suite's query pool...")
	suite.query.Close()
}

func (suite *WebDBTest) SetupTest()    {}
func (suite *WebDBTest) TearDownTest() {}

// Simple test to make sure we can connect to the database, insert data and
// query the results. This is a good general health test to make sure, among
// other things, we can connect to the database.
func (suite *WebDBTest) TestInsertSelect() {
	// enter a dumb test table using `Insert`
	err := suite.entry.Insert("CREATE TABLE test (id INTEGER);")
	defer func() {
		_ = suite.entry.Insert("DROP TABLE test;")
	}()
	if err != nil {
		suite.Fail("unable to create table", err)
	}
	expected := 42
	err = suite.entry.Insert(fmt.Sprintf("INSERT INTO test (id) VALUES (%d);", expected))
	if err != nil {
		suite.Fail("unable to insert into test table", err)
	}

	// get the value using `Select`
	res, err := suite.query.Select("SELECT id FROM test;")
	if err != nil {
		suite.Fail("problem querying test table", err)
	}
	var actual int
	val := res.(pgx.Rows)
	defer val.Close()
	val.Next()
	err = val.Scan(&actual)
	if err != nil {
		suite.Fail("bad reflection", err)
	}

	// verify they're equal
	assert.Equal(suite.T(), expected, actual, "failed simple SQL math")
}

// Insert a bunch of temperature data, get it retreived again
func (suite *WebDBTest) TestInsertSelectTemperatureData() {
	// make a random TempCMap
	size := 100
	temps := generateRandomTempCMap(size)
	logrus.Debug(temps)
}

// make a randomly generated TempCMap
func generateRandomTempCMap(n int) webdb.TempCMap {
	stamps := generateOrderedTimestamps(n)
	temps := make(webdb.TempCMap, n)
	for _, v := range stamps {
		var tempC int
		base := rand.Intn(40)    //nolint:gosec
		neg := rand.Int()%2 == 0 //nolint:gosec
		if neg {
			tempC = base * -1
		} else {
			tempC = base
		}
		temps[v] = tempC
	}
	return temps
}

// get a bunch of ordered timestamps where idx 0 is the oldest and idx -1 is the newest
func generateOrderedTimestamps(num int) []time.Time {
	panic("not implemented, start at generateOrderedTimestamps!")
}
