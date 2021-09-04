package webdb_test

import (
	"fmt"
	"math/rand"
	"sort"
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

const (
	secondsInYear = 60 * 60 * 24 * 365
)

var yearAgo = time.Now().Add(time.Second * -secondsInYear) //nolint:gochecknoglobals

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

func (suite *WebDBTest) SetupTest() {
	// delete all database rows
	err := suite.entry.Insert("DELETE FROM temperature;")
	if err != nil {
		suite.Fail("unable to delete temp tables", err)
	}
	err = suite.entry.Insert("DELETE FROM rain;")
	if err != nil {
		suite.Fail("unable to delete rain tables", err)
	}
}
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
	var expected int32 = 42
	err = suite.entry.Insert(fmt.Sprintf("INSERT INTO test (id) VALUES (%d);", expected))
	if err != nil {
		suite.Fail("unable to insert into test table", err)
	}

	// get the value using `Select`
	res, err := suite.query.Select("SELECT id FROM test;")
	if err != nil {
		suite.Fail("problem querying test table", err)
	}
	actual, err := unwrap(res)

	if err != nil {
		suite.Fail("bad reflection", err)
	}

	// verify they're equal
	assert.Equal(suite.T(), expected, actual, "failed simple SQL math")
}

// Are we actually creating the database from a schema?
func (suite *WebDBTest) TestQueryRealTables() {
	res, err := suite.query.Select("SELECT longname FROM mappings WHERE id=2;")
	if err != nil {
		suite.Fail("failure to SELECT longname FROM mappings", err)
	}
	actual, err := unwrap(res)
	if err != nil {
		suite.Fail("failure to unwrap", err)
	}
	assert.Equal(suite.T(), "soft reset", actual)
}

// Insert a bunch of temperature data, get it retreived again
func (suite *WebDBTest) TestInsertSelectTemperatureData() {
	// make a random TempEntriesC
	size := 100
	expected := generateRandomTempCMap(size)
	for _, entry := range expected {
		err := suite.entry.AddTempCValue(entry.TempC, entry.Timestamp)
		if err != nil {
			suite.Fail("error inserting temperature into database", err)
		}
	}
	since := yearAgo
	actual := suite.query.GetTempDataCSince(since)
	assert.True(suite.T(), len(*actual) == len(expected))
	assert.NotNil(suite.T(), actual)
	for i, v := range *actual {
		// account for subtle rounding errors to go back and forth between postgresql and go
		timeDiff := expected[i].Timestamp.Sub(v.Timestamp)
		assert.True(suite.T(), timeDiff < time.Second)
		assert.Equal(suite.T(), expected[i].TempC, v.TempC, "mismatch on TempEntriesC entry")
	}
}

// Insert data, grab it from a specific range within
func (suite *WebDBTest) TestInsertSelectSpecificTemperatureRange() {
	// make a large chunk of temperature data ordered sequentially by time
	var temps, expected webdb.TempEntriesC
	temp := 0
	start := 4
	end := 6
	beginning := time.Date(2020, time.Month(start), 1, 0, 0, 0, 0, time.UTC)
	finish := time.Date(2020, time.Month(end), 3, 0, 0, 0, 0, time.UTC)
	for i := 1; i < 12; i++ {
		month := time.Month(i)
		timestamp := time.Date(2020, month, 2, 1, 1, 1, 1, time.UTC)
		entry := webdb.TempEntryC{Timestamp: timestamp, TempC: temp}
		temps = append(temps, entry)

		// enter everything into the database
		err := suite.entry.AddTempCValue(temp, timestamp)
		if err != nil {
			suite.Fail("unable to add temp data", err)
		}

		// save a 3-month period to query against
		if int(month) >= start && int(month) <= end {
			expected = append(expected, entry)
		}
		temp++
	}
	// verify the query
	actual := suite.query.GetTempDataCFrom(beginning, finish)
	assert.Equal(suite.T(), len(expected), len(*actual))
	for i, v := range *actual {
		timeDiff := v.Timestamp.Sub(expected[i].Timestamp)
		assert.True(suite.T(), timeDiff < time.Second)
		assert.Equal(suite.T(), expected[i].TempC, v.TempC)
	}
}

/* HELPER FUNCTIONS */

// unwrap a single value
func unwrap(res interface{}) (interface{}, error) {
	var actual interface{}
	val := res.(pgx.Rows)
	defer val.Close()
	val.Next()
	err := val.Scan(&actual)
	if err != nil {
		return nil, err
	}
	return actual, nil
}

// make a randomly generated TempEntriesC
func generateRandomTempCMap(n int) webdb.TempEntriesC {
	var temps []webdb.TempEntryC //nolint:prealloc
	stamps := generateOrderedTimestamps(n)
	for _, stamp := range *stamps {
		var tempC int
		base := rand.Intn(40)    //nolint:gosec
		neg := rand.Int()%2 == 0 //nolint:gosec
		if neg {
			tempC = base * -1
		} else {
			tempC = base
		}
		entry := webdb.TempEntryC{
			Timestamp: stamp,
			TempC:     tempC,
		}
		temps = append(temps, entry)
	}
	return temps
}

// get a bunch of ordered timestamps where idx 0 is the oldest and idx -1 is the newest
func generateOrderedTimestamps(num int) *[]time.Time {
	times := make([]time.Time, num)
	now := time.Now()
	for i := 0; i < num; i++ {
		times[i] = now.Add(time.Second * time.Duration(-rand.Intn(secondsInYear))) //nolint:gosec
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i].Before(times[j])
	})
	return &times
}
