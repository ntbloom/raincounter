package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/spf13/viper"

	"github.com/ntbloom/raincounter/pkg/raincloud/webdb"

	"github.com/sirupsen/logrus"
)

const (
	contentType      = "Content-Type"
	appJSON          = "application/json"
	sensorStatusKey  = "sensor"
	gatewayStatusKey = "gateway"
)

// URLS for handler switch statement
const (
	urlTeapot       = "/v1.0/teapot"
	urlHello        = "/v1.0/hello"
	urlRain         = "/v1.0/rain"
	urlLastRain     = "/v1.0/lastRain"
	urlTemp         = "/v1.0/temp"
	urlLastTemp     = "/v1.0/lastTemp"
	urlSensorStatus = "/v1.0/sensorStatus"
	urlGwStatus     = "/v1.0/gatewayStatus"
)

// restHandler has a connection to the database. Since we're using a read-only
// postgresql connection pool with only GET methods, we don't need a mutex or
// any additional handling. This could change as the application develops.
type restHandler struct {
	db                    webdb.DBQuery
	statusDurationDefault time.Duration
}

// newRestHandler makes a new rest handler with read-only access to the database
func newRestHandler() restHandler {
	logrus.Debug("creating new restHandler")
	return restHandler{
		db:                    webdb.NewPGConnector(),
		statusDurationDefault: viper.GetDuration(configkey.AssetStatusDuration),
	}
}

// close frees any resources needed by the handler
func (handler restHandler) close() {
	logrus.Debug("closing handler struct")
	handler.db.Close()
}

// implement the Handler interface so we can use this as a handler
func (handler restHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// only serve GET requests, all data come in through MQTT
	if r.Method != http.MethodGet {
		logrus.Errorf("attempted illegal request: %s", r.Method)
		return
	}
	logrus.Infof("received request: %s", r.URL.RawQuery)
	switch r.URL.Path {
	case urlHello:
		handler.handleHello(w, r)
	case urlTeapot:
		handler.handleTeapot(w, r)
	case urlLastRain:
		handler.handleLastRain(w, r)
	case urlLastTemp:
		handler.handleLastTemp(w, r)
	case urlSensorStatus:
		handler.handleAssetStatus(sensorStatusKey, w, r)
	case urlGwStatus:
		handler.handleAssetStatus(gatewayStatusKey, w, r)
	case urlTemp:
		handler.handleTemp(w, r)
	case urlRain:
		handler.handleRain(w, r)
	default:
		logrus.Errorf("received unsupported request on `%s`", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}
}

// ParseQuery breaks the restful part of the API into a map
func ParseQuery(raw string) (map[string]interface{}, error) {
	var err error
	if raw == "" {
		err = fmt.Errorf("empty arguments")
		return nil, err
	}
	result := make(map[string]interface{})
	args := strings.Split(raw, "&")
	for _, arg := range args {
		keys := strings.Split(arg, "=")
		if len(keys) < 2 {
			err = fmt.Errorf("illegal REST argument: %s", arg)
			return nil, err
		}
		result[keys[0]] = keys[1]
	}
	return result, nil
}

// handles generic JSON messages. fails if the request does not specify application/json
func (handler restHandler) genericJSONHandler(payload []byte, w http.ResponseWriter, res *http.Request) {
	encoding := res.Header.Get(contentType)
	if encoding != appJSON {
		handler.unsupportedMedia(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		logrus.Error(err)
	}
}

// dateRange is a parsed struct of JSON data with to and from timestamp
type dateRange struct {
	toOk    bool
	fromOk  bool
	totalOk bool
	to      time.Time
	from    time.Time
	total   bool
}

// get args from the rest API
func getToFromTotal(res *http.Request) (*dateRange, error) {
	var err error
	var args map[string]interface{}

	if args, err = ParseQuery(res.URL.RawQuery); err != nil {
		return nil, fmt.Errorf("unparseable arguments")
	}

	_, fromOk := args["from"]
	if !fromOk {
		return nil, err
	}
	var from time.Time
	if from, err = time.Parse(configkey.TimestampFormat, args["from"].(string)); err != nil {
		logrus.Error(err)
		return nil, err
	}
	var to time.Time
	_, toOk := args["to"]
	if toOk {
		to, err = time.Parse(configkey.TimestampFormat, args["to"].(string))
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	total := false
	_, totalOk := args["total"]
	if totalOk {
		if total, err = strconv.ParseBool(args["total"].(string)); err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return &dateRange{
		toOk:    toOk,
		to:      to,
		fromOk:  fromOk,
		from:    from,
		totalOk: totalOk,
		total:   total,
	}, nil
}
