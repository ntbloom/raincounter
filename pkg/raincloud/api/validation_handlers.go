package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

/* HANDLERS FOR GENERIC VALIDATION/TESTING PURPOSES */

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
