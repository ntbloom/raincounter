package database_test

import (
	"os"
	"sync"
	"testing"
	"testing/quick"

	"github.com/ntbloom/raincounter/pkg/common/database"

	"github.com/ntbloom/raincounter/pkg/config"
	"github.com/ntbloom/raincounter/pkg/config/configkey"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

/* FIXTURES */

// reusable configs
func getConfig() {
	config.Configure()
}

// sqliteConnectionFixture makes a reusable Sqlite object
func sqliteConnectionFixture() *database.Sqlite {
	getConfig()
	sqliteFile := viper.GetString(configkey.DatabaseLocalDevFile)
	db, _ := database.NewSqlite(sqliteFile, true)
	return db
}

// Property-based test for creating a bunch of rows and making sure the data get put in
func testRainEntry(db *database.Sqlite, t *testing.T) {
	maxCount := 5
	if testing.Short() {
		logrus.Info("skipping property tests")
		return
	}

	var total int
	test := func(reps uint8) bool {
		count := int(reps)
		for i := 0; i < count; i++ {
			database.MakeRainEntry(db)
		}
		var val int
		if val = database.GetRainEntries(db); val == -1 {
			logrus.Error("gave -1")
			return false
		}
		logrus.Debugf("val=%d, count=%d, total=%d", val, count, total)
		total += count
		return val == total
	}
	if err := quick.Check(test, &quick.Config{
		MaxCount: maxCount,
	}); err != nil {
		t.Error(err)
	}
}

// Tests all the various entries work (except temperature). Also tests concurrent use of postgresql
func testStaticSQLEntries(db *database.Sqlite, t *testing.T) {
	count := 5

	// asynchronously make an entry for each type
	var wg sync.WaitGroup
	wg.Add(5 * count)
	type addFunction func(db database.DBWrapper)
	checkAdd := func(callable addFunction, arg database.DBWrapper) {
		defer wg.Done()
		callable(arg)
	}
	for i := 0; i < count; i++ {
		go checkAdd(database.MakeRainEntry, db)
		go checkAdd(database.MakeSoftResetEntry, db)
		go checkAdd(database.MakeHardResetEntry, db)
		go checkAdd(database.MakePauseEntry, db)
		go checkAdd(database.MakeUnpauseEntry, db)
	}
	// wait for entries to finish
	wg.Wait()

	// verify counts
	wg.Add(5)
	type getFunction func(db database.DBWrapper) int
	checkGet := func(callable getFunction, arg database.DBWrapper) {
		defer wg.Done()
		tally := callable(arg)
		if tally != count {
			t.Fail()
		}
	}
	go checkGet(database.GetRainEntries, db)
	go checkGet(database.GetSoftResetEntries, db)
	go checkGet(database.GetHardResetEntries, db)
	go checkGet(database.GetPauseEntries, db)
	go checkGet(database.GetUnpauseEntries, db)
	wg.Wait()
}

// tests that we can enter temperature
func testTemperatureEntries(db *database.Sqlite, t *testing.T) {
	vals := []int{-100, -25, -15, -1, 0, 1, 2, 20, 24, 100}
	for _, expected := range vals {
		database.MakeTemperatureEntry(db, expected)
		if actual := database.GetLastTemperatureEntry(db); expected != actual {
			logrus.Errorf("expected=%d, actual=%d", expected, actual)
			t.Fail()
		}
	}
}

/* TESTS */

/* Starting with Sqlite, make sure the schema and file manipulation are enforced properly */

// create and destroy sqlite file 5 times, get Sqlite Sqlite struct
func TestSqliteDataPrep(t *testing.T) {
	getConfig()
	sqliteFile := viper.GetString(configkey.DatabaseLocalDevFile)

	// clean up when finished
	defer func() { _ = os.Remove(sqliteFile) }()

	// create and destroy 5 times
	for i := 0; i < 5; i++ {
		db, err := database.NewSqlite(sqliteFile, true)
		if err != nil || db == nil {
			logrus.Error("problem instantiating NewSqlite struct")
			t.Error(err)
		}
		_, err = os.Stat(sqliteFile)
		if err != nil {
			logrus.Error("sqlite file doesn't exist")
			t.Error(err)
		}
	}
}

func TestSqliteForeignKeysAreImplemented(t *testing.T) {
	db := sqliteConnectionFixture()
	if foreignKeys := db.ForeignKeysAreImplemented(); !foreignKeys {
		logrus.Error("sqlite is not enforcing foreign_key constraints")
		t.Fail()
	}
}

func TestSqliteRainEntry(t *testing.T) {
	db := sqliteConnectionFixture()
	testRainEntry(db, t)
}

func TestSqliteStaticSQLEntries(t *testing.T) {
	db := sqliteConnectionFixture()
	testStaticSQLEntries(db, t)
}

func TestSqliteTemperatureEntries(t *testing.T) {
	db := sqliteConnectionFixture()
	testTemperatureEntries(db, t)
}
