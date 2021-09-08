package webdb_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/ntbloom/raincounter/pkg/gateway/tlv"

	"github.com/stretchr/testify/assert"

	"github.com/sirupsen/logrus"

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
	rainAmt float64
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
	suite.rainAmt = viper.GetFloat64(configkey.SensorRainMm)
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
	for _, sql := range []string{
		"DELETE FROM temperature;",
		"DELETE FROM rain;",
		"DELETE FROM event_log;",
		"DELETE FROM status_log;",
	} {
		err := suite.entry.Insert(sql)
		if err != nil {
			suite.Fail("can't delete table rows", err)
		}
	}
}
func (suite *WebDBTest) TearDownTest() {}

/* GENERIC PACKAGE-LEVEL TESTS */

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

/* TEMPERATURE INSERT/QUERY TESTS */

// Insert a bunch of temperature data, get it retreived again
func (suite *WebDBTest) TestInsertSelectTemperatureData() {
	// make a random TempEntriesC
	size := 100
	expected := generateRandomTempEntriesC(size)
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

// Insert temperature data, grab it from within a specific range
func (suite *WebDBTest) TestInsertSelectSpecificTemperatureRange() {
	// make a large chunk of temperature data ordered sequentially by time
	var expected webdb.TempEntriesC
	temp := 0

	// isolate a subset of the whole year
	start := 4
	end := 6
	beginning := time.Date(2020, time.Month(start), 1, 0, 0, 0, 0, time.UTC)
	finish := time.Date(2020, time.Month(end), 3, 0, 0, 0, 0, time.UTC)

	// enter data for the whole year, set aside ones that fit in the subset
	for i := 1; i < 12; i++ {
		month := time.Month(i)
		timestamp := time.Date(2020, month, 2, 1, 1, 1, 1, time.UTC)
		entry := webdb.TempEntryC{Timestamp: timestamp, TempC: temp}

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

// Can we get the last temperature value
func (suite *WebDBTest) TestGetLastTempC() {
	randomData := generateRandomTempEntriesC(100)
	var maxDate time.Time
	var maxTemp int
	for i, v := range randomData {
		stamp := v.Timestamp
		temp := v.TempC
		if i == 0 {
			maxDate = stamp
			maxTemp = temp
		}
		if stamp.After(maxDate) {
			maxDate = stamp
			maxTemp = temp
		}
		err := suite.entry.AddTempCValue(temp, stamp)
		if err != nil {
			suite.Fail("error inserting temp data", err)
		}
	}
	actual := suite.query.GetLastTempC()
	assert.Equal(suite.T(), maxTemp, actual)
}

// make a randomly generated TempEntriesC
func generateRandomTempEntriesC(n int) webdb.TempEntriesC {
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

/* RAIN INSERT/QUERY TESTS */

// enter 100 rows and retries them, don't worry about date range
func (suite *WebDBTest) TestEnterAndRetrieveRainAllData() {
	// generate random array of data
	data := generateRandomRainEntriesMM(100)
	var expTotalRain float64 = 0.0
	for _, entry := range data {
		err := suite.entry.AddRainMMEvent(entry.Millimeters, entry.Timestamp)
		expTotalRain += entry.Millimeters
		if err != nil {
			suite.Fail("failed to add rain amount", err)
		}
	}
	actual := suite.query.GetRainMMSince(yearAgo)
	assert.NotNil(suite.T(), actual)
	var actTotalRain float64 = 0.0
	for i, v := range *actual {
		actTotalRain += v.Millimeters
		timeDiff := data[i].Timestamp.Sub(v.Timestamp)
		assert.True(suite.T(), timeDiff < time.Second)
	}
	assert.Equal(suite.T(), len(data), len(*actual), "length of actual and expected data are not equal")
	assert.Equal(suite.T(), expTotalRain, actTotalRain, "actual and total rain are not equal")
}

// test all the rain values on a selected range
func (suite *WebDBTest) TestEnterAndRetrieveRainDataWithinRange() {
	amt := viper.GetFloat64(configkey.SensorRainMm)
	var expected webdb.RainEntriesMm
	var expTotal float64 = 0.0

	// set a time range to query against
	start := 4
	end := 6
	beginning := time.Date(2020, time.Month(start), 1, 0, 0, 0, 0, time.UTC)
	finish := time.Date(2020, time.Month(end), 3, 0, 0, 0, 0, time.UTC)

	// add data for the whole year, set aside the subrange
	for i := 1; i < 12; i++ {
		month := time.Month(i)
		timestamp := time.Date(2020, month, 2, 1, 1, 1, 1, time.UTC)
		entry := webdb.RainEntryMm{Timestamp: timestamp, Millimeters: amt}

		// enter everything into the database
		err := suite.entry.AddRainMMEvent(amt, timestamp)
		if err != nil {
			suite.Fail("unable to add rain data", err)
		}

		// save a 3-month period to query against
		if int(month) >= start && int(month) <= end {
			expected = append(expected, entry)
			expTotal += amt
		}
	}
	actual := suite.query.GetRainMMFrom(beginning, finish)
	assert.NotNil(suite.T(), actual)
	var actTotal float64 = 0.0
	for i, v := range *actual {
		assert.Equal(suite.T(), v.Millimeters, amt, "amount was entered incorrectly")
		timeDiff := v.Timestamp.Sub(expected[i].Timestamp)
		assert.True(suite.T(), timeDiff < time.Second, "timestamps do not match")
		actTotal += amt
	}
	assert.Equal(suite.T(), expTotal, actTotal, "total amounts are unequal")

	// also verify the tallying function
	queriedTotal := suite.query.TotalRainMMFrom(beginning, finish)
	assert.Equal(suite.T(), expTotal, queriedTotal, "rain tallying function is incorrect")
}

// Insert some rain values, get timestamp from the last entry entered
func (suite *WebDBTest) TestGetLastRainTime() {
	// make 2 timestamps, enter rain event for it
	amt := viper.GetFloat64(configkey.SensorRainMm)
	twoHoursAgo := time.Now().Add(time.Hour * -2)
	oneHourAgo := time.Now().Add(time.Hour * -1)
	for _, stamp := range []time.Time{twoHoursAgo, oneHourAgo} {
		err := suite.entry.AddRainMMEvent(amt, stamp)
		if err != nil {
			suite.Fail("failed to enter value", err)
		}
	}
	//
	lastRainTime := suite.query.GetLastRainTime()
	timeDiff := lastRainTime.Sub(oneHourAgo)
	assert.True(suite.T(), timeDiff < time.Second)
}

// generate a random RainEntriesMM struct
func generateRandomRainEntriesMM(n int) webdb.RainEntriesMm {
	stamps := generateOrderedTimestamps(n)
	amt := viper.GetFloat64(configkey.SensorRainMm)
	var rain webdb.RainEntriesMm
	for _, stamp := range *stamps {
		entry := webdb.RainEntryMm{
			Timestamp: stamp, Millimeters: amt,
		}
		rain = append(rain, entry)
	}
	return rain
}

// make sure we don't error on event/status messages
func (suite *WebDBTest) TestEventAndStatusMessagesDontError() {
	for _, asset := range []int{configkey.SensorStatus, configkey.GatewayStatus} {
		err := suite.entry.AddStatusUpdate(asset, time.Now())
		if err != nil {
			suite.Fail("unable to add status message", err)
		}
	}
	for _, tag := range []int{tlv.SoftReset, tlv.HardReset, tlv.Pause, tlv.Unpause} {
		err := suite.entry.AddTagValue(tag, 1, time.Now())
		if err != nil {
			suite.Fail("unable to add tag", err)
		}
	}
	// do a few quick and dirty sql queries just to make sure something made it into the database

	// status page
	statusQuery := `SELECT sum(asset) FROM status_log;` // should be 1(sensor) + 2(gateway), so 3
	val, err := suite.query.Select(statusQuery)
	if err != nil {
		suite.Fail("unable to query status_log table", err)
	}
	statuses, err := unwrap(val)
	if err != nil {
		suite.Fail("unable to unwrap statuses", err)
	}
	assert.Equal(suite.T(), int64(3), statuses)

	tagQuery := `SELECT sum(value) FROM event_log;` // should be 4, one for each tlv tag
	val, err = suite.query.Select(tagQuery)
	if err != nil {
		suite.Fail("unable to query event_log table", err)
	}
	tags, err := unwrap(val)
	if err != nil {
		suite.Fail("error unwrapping event_log", err)
	}
	assert.Equal(suite.T(), int64(4), tags)
}

// make sure we handle the database rows being empty
func (suite *WebDBTest) TestEmptyResultsDontError() {
	// test assumes all rows are empty

	rainSince := suite.query.GetRainMMSince(time.Now())
	assert.Zero(suite.T(), len(*rainSince), "expected an empty struct")

	rainBetween := suite.query.TotalRainMMFrom(time.Now(), time.Now())
	assert.Zero(suite.T(), rainBetween, "function should return 0.0 when no matches")

	tempSince := suite.query.GetTempDataCSince(time.Now())
	assert.Zero(suite.T(), len(*tempSince), "expected an empty struct")

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
