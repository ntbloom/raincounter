package api

import (
	"net/http"
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
