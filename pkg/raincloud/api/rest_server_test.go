package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ntbloom/raincounter/pkg/raincloud/webdb"

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
	rest  *api.RestServer
	entry webdb.DBEntry
	url   string
}

func TestApi(t *testing.T) {
	test := new(RestTest)
	suite.Run(t, test)
}

func (suite *RestTest) SetupSuite() {
	config.Configure()

	// add a query connector
	suite.entry = webdb.NewPGConnector()

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

// read
func (suite *RestTest) toJSON(resp *http.Response, passedErr error) ([]byte, int) {
	status := resp.StatusCode
	if passedErr != nil {
		suite.Fail("error getting response", passedErr)
		return nil, status
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			suite.Fail("failure to close body", err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.Fail("error reading response body", err)
	}
	return body, status
}

func (suite *RestTest) connectToServer() bool {
	var resp *http.Response
	var err error
	for i := 0; i < 20; i++ {
		resp, err = suite.getEndpoint("/teapot")
		if resp != nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.Fail("error reading empty message", err)
	}
	assert.Equal(suite.T(), "", string(body), "should not be returning any payload")
	assert.Equal(suite.T(), http.StatusUnsupportedMediaType, resp.StatusCode)
}

// make sure we can connect to the API
func (suite *RestTest) TestHello() {
	body, status := suite.toJSON(suite.getEndpoint("/hello"))

	expected := `{"hello":"world"}`
	actual := string(body)

	assert.Equal(suite.T(), http.StatusOK, status)
	assert.Equal(suite.T(), expected, actual)
}

// get the last rain value as a timestamp
func (suite *RestTest) TestGetLastRain() {
	rain, status := suite.toJSON(suite.getEndpoint("/lastRain"))
	var actual map[string]time.Time
	err := json.Unmarshal(rain, &actual)

	assert.Equal(suite.T(), http.StatusOK, status)
	assert.NotNil(suite.T(), actual)
	assert.Equal(suite.T(), 1, len(actual), "should only have 1 result")
	assert.Nil(suite.T(), err)
}

// get the last temperature value as an integer
func (suite *RestTest) TestGetLastTempC() {
	temp, status := suite.toJSON(suite.getEndpoint("/lastTemp"))
	var actual map[string]int
	err := json.Unmarshal(temp, &actual)

	assert.Equal(suite.T(), http.StatusOK, status)
	assert.NotNil(suite.T(), actual)
	assert.Equal(suite.T(), 1, len(actual), "should only have 1 result")
	assert.Nil(suite.T(), err)
}

// test the rest API argument parser
func (suite *RestTest) TestParseQuery() {
	// we can afford to be slapdash and only support the patterns we are actually coding
	args := map[string]map[string]interface{}{
		"since=300": {"since": "300"},
		"from=2021-09-23T01:22:18+00:00&to=2021-09-23T01:22:18+00:00": {
			"from": "2021-09-23T01:22:18+00:00",
			"to":   "2021-09-23T01:22:18+00:00",
		},
	}
	for k, v := range args {
		expected, err := api.ParseQuery(k)
		if err != nil {
			suite.Fail("error parsing query", err)
		}
		assert.Equal(suite.T(), v, expected)
	}
}

// test receiving OK for sensor and gateway status, first with a false and then a true
func (suite *RestTest) TestGetStatus() {
	statusTest := func(endpoint string, activeKey string, statusNum int) {
		// we expect there to not be anything at the beginning since the dummy data are old
		var actual map[string]interface{}
		var err error

		beforeStatus, status := suite.toJSON(suite.getEndpoint(endpoint))
		err = json.Unmarshal(beforeStatus, &actual)
		inactive := actual[activeKey].(bool)

		assert.Equal(suite.T(), http.StatusOK, status)
		assert.False(suite.T(), inactive, "shouldn't be an entry yet")
		assert.Nil(suite.T(), err)

		// add a more updated entry, remove it when we're done
		if err = suite.entry.AddStatusUpdate(statusNum, time.Now()); err != nil {
			suite.Fail("problem entering status message", err)
		}
		defer func() {
			cmd := `DELETE FROM status_log WHERE id=(SELECT id FROM status_log ORDER BY gw_timestamp DESC LIMIT 1);`
			if err = suite.entry.Insert(cmd); err != nil {
				logrus.Warning("did not erase last sensor entry")
			}
		}()

		afterStatus, status := suite.toJSON(suite.getEndpoint(endpoint))
		err = json.Unmarshal(afterStatus, &actual)
		active := actual[activeKey].(bool)

		assert.Equal(suite.T(), http.StatusOK, status)
		assert.True(suite.T(), active, "should be picked up")
		assert.Nil(suite.T(), err)
	}

	// run the test for the gateway and sensor test
	for _, v := range map[string]map[string]interface{}{
		"sensor":  {"endpoint": "/sensorStatus?since=300", "activeKey": "sensor_active", "status": configkey.SensorStatus},
		"gateway": {"endpoint": "/gatewayStatus?since=300", "activeKey": "gateway_active", "status": configkey.GatewayStatus},
	} {
		endpoint := v["endpoint"].(string)
		activeKey := v["activeKey"].(string)
		status := v["status"].(int)
		statusTest(endpoint, activeKey, status)
	}
}

func (suite *RestTest) TestGetTemperatureData() {
	sampleSince := "/temp?since=2021-09-23T01:47:30+00:00"
	//sampleFrom := "/temp?from=2021-05-23T01:22:18+00:00&to=2021-09-23T01:22:18+00:00"

	var tempResults []map[string]interface{}

	since, statusSince := suite.toJSON(suite.getEndpoint(sampleSince))
	if err := json.Unmarshal(since, &tempResults); err != nil {
		suite.Fail("unable to unmarshal json", err)
	}
	assert.Equal(suite.T(), http.StatusOK, statusSince)
	assert.NotNil(suite.T(), since)
	assert.Greater(suite.T(), len(since), 0, "results are empty")
}

/* NEED TO WRITE ENDPOINTS FOR THE FOLLOWING ENDPOINTS */

//TotalRainMMSince(since time.Time) (float64, error)
//TotalRainMMFrom(from time.Time, to time.Time) (float64, error)
//GetRainMMSince(since time.Time) (*RainEntriesMm, error)
//GetRainMMFrom(from time.Time, to time.Time) (*RainEntriesMm, error)
//GetTempDataCSince(since time.Time) (*TempEntriesC, error)
//GetTempDataCFrom(from time.Time, to time.Time) (*TempEntriesC, error)
//GetEventMessagesSince(tag int, since time.Time) (*EventEntries, error)
//GetEventMessagesFrom(tag int, from time.Time, to time.Time) (*EventEntries, error)
