package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

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
	if payload, err := json.Marshal(map[string]string{"hello": "world"}); err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		genericJSONHandler(payload, w, res)
	}
}

func genericJSONHandler(payload []byte, w http.ResponseWriter, res *http.Request) {
	encoding := res.Header.Get(contentType)
	if encoding != appJson {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		logrus.Error(err)
	}
}
