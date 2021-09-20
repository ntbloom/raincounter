package api_test

import (
	"context"
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

/* TESTING FIXTURES */

type RestTest struct {
	suite.Suite
	rest *api.RestServer
	url  string
}

func TestApi(t *testing.T) {
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

func (suite *RestTest) TearDownSuite() {
	logrus.Info("stopping the rest API from the test suite")
	suite.rest.Stop()
}

/* HELPER METHODS */

// just get the response from a GET, fail if there's an error
func (suite *RestTest) getEndpoint(endpoint string) (*http.Response, error) {
	var err error
	var req *http.Request
	var resp *http.Response

	url := fmt.Sprintf("%s%s", suite.url, endpoint)
	var headers = map[string]string{
		"content-type": "application/json",
	}

	if req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil); err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err = http.DefaultClient.Do(req)
	return resp, err
}

func (suite *RestTest) connectToServer() bool {
	var resp *http.Response
	var err error
	for i := 0; i < 20; i++ {
		resp, err = suite.getEndpoint("/teapot")
		if resp != nil {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			suite.Fail("not a teapot!", err)
		}
	}()
	assert.Nil(suite.T(), err, fmt.Sprintf("error retreiving teapot: %s", err))
	assert.Equal(suite.T(), http.StatusTeapot, resp.StatusCode)
	return true
}

/* TESTS */

func (suite *RestTest) TestTeapot() {
	resp, err := suite.getEndpoint("/teapot") //nolint:bodyclose
	if err != nil {
		suite.Fail("problem getting teapot", err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			suite.Fail("error closing body", err)
		}
	}()
	assert.Equal(suite.T(), http.StatusTeapot, resp.StatusCode)
}

func (suite *RestTest) TestHello() {
	resp, err := suite.getEndpoint("/hello") // nolint:bodyclose
	if err != nil {
		suite.Fail("problem getting hello", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			suite.Fail("failed to close body", err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.Fail("error reading response", err)
	}
	expected := "{\"hello\":\"world\"}"
	actual := string(body)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.Equal(suite.T(), expected, actual)
}

// make sure we get a bad response if we forget to set application/json content-type header
func (suite *RestTest) TestNoJsonHeaders() {
	var resp *http.Response
	var err error
	url := fmt.Sprintf("%s%s", suite.url, "/hello")
	resp, err = http.Get(url) //nolint
	if resp != nil {
		defer func() {
			if err = resp.Body.Close(); err != nil {
				suite.Fail("error closing hello", err)
			}
		}()
	}
	assert.Equal(suite.T(), http.StatusUnsupportedMediaType, resp.StatusCode)
}
