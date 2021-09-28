package api

import (
	"encoding/json"
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

/* HELPER METHODS */

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
	toOk   bool
	fromOk bool
	to     time.Time
	from   time.Time
}

func getToFrom(res *http.Request) (*dateRange, error) {
	args, err := ParseQuery(res.URL.RawQuery)
	if err != nil {
		logrus.Error(err)
		return nil, err
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
	return &dateRange{
		toOk:   toOk,
		to:     to,
		fromOk: fromOk,
		from:   from,
	}, nil
}

// ParseQuery breaks the restful part of the API into a map
func ParseQuery(raw string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	args := strings.Split(raw, "&")
	for _, arg := range args {
		keys := strings.Split(arg, "=")
		result[keys[0]] = keys[1]
	}
	return result, nil
}

/* GENERIC AND TEST HANDLERS */

// return teapot messages as bellweather for general server and for bootstrapping
// may be able to delete this later as the API is developed
func (handler restHandler) handleTeapot(w http.ResponseWriter, _ *http.Request) {
	var payload []byte
	var err error

	if payload, err = json.Marshal(map[string]string{"hello": "teapot"}); err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusTeapot)
		if _, err = w.Write(payload); err != nil {
			logrus.Error(err)
		}
	}
}

// template for json payload messages
func (handler restHandler) handleHello(w http.ResponseWriter, res *http.Request) {
	payload, err := json.Marshal(map[string]string{"hello": "world"})
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	handler.genericJSONHandler(payload, w, res)
}

/* PRODUCTION ENDPOINT HANDLERS */

// handle requests for the last rain
func (handler restHandler) handleLastRain(w http.ResponseWriter, res *http.Request) {
	payload, err := handler.db.GetLastRainTime()
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(map[string]time.Time{"timestamp": payload})
	if err != nil {
		logrus.Error(err)
	}
	handler.genericJSONHandler(resp, w, res)
}

// handle requests for the last temp
func (handler restHandler) handleLastTemp(w http.ResponseWriter, res *http.Request) {
	payload, err := handler.db.GetLastTempC()
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(map[string]int{"last_temp_c": payload})
	if err != nil {
		logrus.Error(err)
	}
	handler.genericJSONHandler(resp, w, res)
}

// handle requests for sensor status
func (handler restHandler) handleAssetStatus(asset string, w http.ResponseWriter, res *http.Request) {
	var dbQuery func(duration time.Duration) (bool, error)
	var responseKey string
	switch asset {
	case sensorStatusKey:
		dbQuery = handler.db.IsSensorUp
		responseKey = "sensor_active"
	case gatewayStatusKey:
		dbQuery = handler.db.IsGatewayUp
		responseKey = "gateway_active"
	}

	raw := res.URL.RawQuery
	args, err := ParseQuery(raw)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	since := handler.statusDurationDefault
	_, ok := args["since"]
	if ok {
		asNum, err := strconv.Atoi(args["since"].(string))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		since = time.Second * time.Duration(asNum)
	}
	isUp, err := dbQuery(since)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(map[string]interface{}{responseKey: isUp})
	if err != nil {
		logrus.Error(err)
	}
	handler.genericJSONHandler(resp, w, res)
}

// handle request for temperature data
func (handler restHandler) handleTemp(w http.ResponseWriter, res *http.Request) {
	var dates *dateRange
	var err error
	var entries *webdb.TempEntriesC
	var resp []byte

	if dates, err = getToFrom(res); err != nil {
		handler.badRequest(w, err)
		return
	}
	if dates.toOk {
		if entries, err = handler.db.GetTempDataCFrom(dates.from, dates.to); err != nil {
			handler.internalServiceError(w, err)
			return
		}
	} else {
		if entries, err = handler.db.GetTempDataCSince(dates.from); err != nil {
			handler.internalServiceError(w, err)
			return
		}
	}
	if resp, err = json.Marshal(entries); err != nil {
		handler.internalServiceError(w, err)
		return
	}
	handler.genericJSONHandler(resp, w, res)
}

// handle request for rain data
func (handler restHandler) handleRain(w http.ResponseWriter, res *http.Request) {
	var dates *dateRange
	var err error
	var entries *webdb.RainEntriesMm
	var resp []byte

	if dates, err = getToFrom(res); err != nil {
		handler.badRequest(w, err)
		return
	}
	if dates.toOk {
		if entries, err = handler.db.GetRainMMFrom(dates.from, dates.to); err != nil {
			handler.internalServiceError(w, err)
		}
	} else {
		if entries, err = handler.db.GetRainMMSince(dates.from); err != nil {
			handler.internalServiceError(w, err)
		}
	}
	if resp, err = json.Marshal(entries); err != nil {
		handler.internalServiceError(w, err)
		return
	}
	handler.genericJSONHandler(resp, w, res)
}
