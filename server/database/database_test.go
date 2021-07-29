package database_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	"github.com/ntbloom/raincloud/pkg/database"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// VARIOUS HELPER FUNCTIONS

const (
	podman = "/usr/bin/podman"
)

// run a generic command
func bashCommand(executable string, args []string) {
	cmd := exec.Command(executable, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		logrus.Fatal(err)
	}
}

// launch the database
func launchDatabase(container string, database string) {
	args := []string{
		"run",
		"--name",
		container,
		"--rm",
		"-d",
		"-e",
		"POSTGRES_PASSWORD=password",
		"-e",
		fmt.Sprintf("POSTGRES_DB=%s", database),
		"-p",
		"5432:5432",
		"docker.io/library/postgres",
	}
	bashCommand(podman, args)
}

// kill the database
func killDatabase(name string) {
	args := []string{"kill", name}
	bashCommand(podman, args)
}

// FIXTURES

type DatabaseTestSuite struct {
	suite.Suite
	containerName string
	dbName        string
	conn          *database.PgConnector
}

// run once at startup
func (suite *DatabaseTestSuite) SetupSuite() {
	launchDatabase(suite.containerName, suite.dbName)
}

// run once at the end
func (suite *DatabaseTestSuite) TearDownSuite() {
	killDatabase(suite.containerName)
}

// run at end of each test
func (suite *DatabaseTestSuite) TearDownTest() {
}

// run at start of each test
func (suite *DatabaseTestSuite) SetupTest() {
}

// ACTUAL TESTS

func (suite *DatabaseTestSuite) TestBasicConnection() {
	err := suite.conn.MakeContact()
	assert.Nil(suite.T(), err)
}

func TestPostgresSuite(t *testing.T) {
	dbTest := new(DatabaseTestSuite)
	dbTest.containerName = "test_postgres_container"
	dbTest.dbName = "raincloud_test"

	const url = "postgresql://postgres:password@localhost:5432/raincloud_test"
	dbTest.conn = database.NewDatabase(dbTest.dbName, url)
	suite.Run(t, dbTest)
}
