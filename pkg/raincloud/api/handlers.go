package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	contentType = "Content-Type"
	appJSON     = "application/json"
)

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
func handleTeapot(w http.ResponseWriter, res *http.Request) {
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
func handleHello(w http.ResponseWriter, res *http.Request) {
	payload, err := json.Marshal(map[string]string{"hello": "world"})
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	genericJSONHandler(payload, w, res)
}

/* PRODUCTION ENDPOINT HANDLERS */

// handle rain requests
func handleRain(w http.ResponseWriter, res *http.Request) {
	panic("implement me!")
}
