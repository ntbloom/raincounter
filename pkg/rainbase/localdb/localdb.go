package localdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ntbloom/raincounter/pkg/common/database"

	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite" // Driver for localdb
)

type LocalDB struct {
	lite *database.Sqlite
}

func NewLocalDB(fulPath string, clobber bool) (*LocalDB, error) {
	lite, err := database.NewSqlite(fulPath, clobber, localDbSchema)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &LocalDB{lite}, nil
}

func (db *LocalDB) MakeSchema() (sql.Result, error) {
	return db.lite.MakeSchema(localDbSchema)
}

func (db *LocalDB) EnterData(cmd string) (sql.Result, error) {
	return db.lite.EnterData(cmd)
}

func (db *LocalDB) AddIntRecord(tag, value int) (sql.Result, error) {
	timestamp := time.Now().Format(time.RFC3339)
	cmd := fmt.Sprintf("INSERT INTO log (tag, value, timestamp) VALUES (%d, %d, \"%s\");", tag, value, timestamp)
	return db.lite.EnterData(cmd)
}

func (db *LocalDB) AddFloatRecord(tag int, value float64) (sql.Result, error) {
	panic("float not implemented on the rainbase")
}

func (db *LocalDB) Tally(tag int) int {
	query := fmt.Sprintf("SELECT COUNT(*) FROM log WHERE tag = %d;", tag)
	return db.GetSingleInt(query)
}

func (db *LocalDB) GetLastRecord(tag int) int {
	cmd := fmt.Sprintf(`SELECT value FROM log WHERE tag = %d ORDER BY id DESC LIMIT 1;`, tag)
	return db.GetSingleInt(cmd)
}

func (db *LocalDB) GetSingleInt(query string) int {
	var rows *sql.Rows
	var err error

	c, _ := db.lite.Connect() // don't handle the error, just return -1
	defer c.Disconnect()

	if rows, err = c.Conn.QueryContext(context.Background(), query); err != nil {
		return -1
	}
	closed := func() {
		if err = rows.Close(); err != nil {
			logrus.Error(err)
		}
	}
	defer closed()
	results := make([]int, 0)
	for rows.Next() {
		var val int
		if err = rows.Scan(&val); err != nil {
			logrus.Error(err)
			return -1
		}
		results = append(results, val)
	}

	return results[0]
}

// ForeignKeysAreImplemented tests function to ensure foreign key implementation
func (db *LocalDB) ForeignKeysAreImplemented() bool {
	illegal := `INSERT INTO log (tag, value, timestamp) VALUES (99999,1,"timestamp");`
	res, err := db.lite.EnterData(illegal)
	return res == nil && err != nil
}
