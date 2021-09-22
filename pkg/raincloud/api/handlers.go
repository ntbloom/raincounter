package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ntbloom/raincounter/pkg/raincloud/webdb"

	"github.com/sirupsen/logrus"
)

const (
	contentType = "Content-Type"
	appJSON     = "application/json"
)

type restHandler struct {
	db webdb.DBQuery
}

func newRestHandler() restHandler {
	logrus.Debug("creating new restHandler")
	var query webdb.DBQuery
	db := webdb.NewPGConnector()
	query = db
	return restHandler{db: query}
}

func (rest restHandler) close() {
	logrus.Info("closing handler struct")
	rest.db.Close()
}

func (rest restHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logrus.Errorf("attempted illegal request: %s", r.Method)
		return
	}
	switch r.URL.Path {
	case hello:
		rest.handleHello(w, r)
	case teapot:
		rest.handleTeapot(w, r)
	case lastRain:
		rest.handleLastRain(w, r)
	}
}

/* GENERIC AND TEST HANDLERS */

// handles generic JSON messages. fails if the request does not specify application/json
func genericJSONHandler(payload []byte, w http.ResponseWriter, res *http.Request) {
	encoding := res.Header.Get(contentType)
	if encoding != appJSON {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		logrus.Error(err)
	}
}

// return teapot messages as bellweather for general server and for bootstrapping
// may be able to delete this later as the API is developed
func (rest restHandler) handleTeapot(w http.ResponseWriter, res *http.Request) {
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
func (rest restHandler) handleHello(w http.ResponseWriter, res *http.Request) {
	payload, err := json.Marshal(map[string]string{"hello": "world"})
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	genericJSONHandler(payload, w, res)
}

/* PRODUCTION ENDPOINT HANDLERS */

// handle requests for the last rain
func (rest restHandler) handleLastRain(w http.ResponseWriter, res *http.Request) {
	payload, err := rest.db.GetLastRainTime()
	if err != nil {
		logrus.Error(err)
		return
	}
	resp, err := json.Marshal(map[string]time.Time{"timestamp": payload})
	if err != nil {
		logrus.Error(err)
	}
	genericJSONHandler(resp, w, res)
}

// handle rain requests
func handleRain(w http.ResponseWriter, res *http.Request) {
	panic("implement me!")
}
