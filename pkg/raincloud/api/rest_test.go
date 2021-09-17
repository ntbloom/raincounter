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

	// wait for the server to be up
	assert.True(suite.T(), suite.connectToServer(), "unable to connect to server")
}

func (suite *RestTest) connectToServer() bool {
	var resp *http.Response
	var err error
	for i := 0; i < 5; i++ {
		url := fmt.Sprintf("%s%s", suite.url, "/teapot")
		resp, err = http.Get(url) //nolint
		if resp != nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			suite.Fail("not a teapot!", err)
		}
	}()
	assert.Nil(suite.T(), err, fmt.Sprintf("error retreiving teapot: %s", err))
	assert.Equal(suite.T(), resp.StatusCode, http.StatusTeapot)
	return true
}

func (suite *RestTest) TearDownSuite() {
	logrus.Info("stopping the rest API from the test suite")
	suite.rest.Stop()
}
func (suite *RestTest) SetupTest()    {}
func (suite *RestTest) TearDownTest() {}

func (suite *RestTest) TestHelloWorld() {
	resp := suite.callEndpoint("/hello") //nolint:bodyclose
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			suite.Fail("error closing body", err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.Fail("error reading response", err)
	}
	message := string(body)

	assert.Equal(suite.T(), 200, resp.StatusCode)
	assert.Equal(suite.T(), "Hello, world!", message)
}

// just get the response from a , fail if there's an error
func (suite *RestTest) callEndpoint(endpoint string) *http.Response {
	url := fmt.Sprintf("%s%s", suite.url, endpoint)
	resp, err := http.Get(url) //nolint
	if err != nil {
		suite.Fail(fmt.Sprintf("failure to call %s", url), err)
	}
	return resp
}
