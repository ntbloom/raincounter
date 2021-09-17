package api_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"

	"github.com/ntbloom/raincounter/pkg/config"
	"github.com/ntbloom/raincounter/pkg/raincloud/api"
	"github.com/stretchr/testify/suite"
)

type RestTest struct {
	suite.Suite
	rest *api.RestServer
	url  string
}

func TestReceiver(t *testing.T) {
	test := new(RestTest)
	suite.Run(t, test)
}

func (suite *RestTest) SetupSuite() {
	config.Configure()

	// launch the rest API
	rest, err := api.NewRestServer()
	if err != nil {
		suite.Fail("error making rest server", err)
	}
	suite.rest = rest
	go suite.rest.Run()

	// get a base API to query against
	scheme := viper.GetString(configkey.RestScheme)
	baseurl := viper.GetString(configkey.RestIP)
	port := viper.GetString(configkey.RestPort)
	version := viper.Get(configkey.RestVersion)
	suite.url = fmt.Sprintf("%s://%s:%s/%s", scheme, baseurl, port, version)
}

func (suite *RestTest) TearDownSuite() {
	logrus.Info("stopping the rest API from the test suite")
	suite.rest.Stop()
}
func (suite *RestTest) SetupTest()    {}
func (suite *RestTest) TearDownTest() {}

func (suite *RestTest) TestHelloWorld() {
	time.Sleep(time.Millisecond * 500)
	url := fmt.Sprintf("%s/hello", suite.url)
	logrus.Debugf("url=%s", url)
	resp, err := http.Get(url)
	if err != nil {
		suite.Fail("error getting hello world", err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			suite.Fail("error closing body", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.Fail("error reading response", err)
	}
	message := string(body)
	logrus.Error(message)

	assert.Equal(suite.T(), 200, resp.StatusCode)
	assert.Equal(suite.T(), "Hello, world!", message)
}
