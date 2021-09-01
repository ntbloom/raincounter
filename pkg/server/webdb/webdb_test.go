package webdb_test

import (
	"fmt"
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

// Simple test to make sure we can connect to the database, insert data and query the results
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
