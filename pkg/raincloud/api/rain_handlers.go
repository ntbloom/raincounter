package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ntbloom/raincounter/pkg/raincloud/webdb"
	"github.com/sirupsen/logrus"
)

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

// handle request for rain data
func (handler restHandler) handleRain(w http.ResponseWriter, res *http.Request) {
	var dates *dateRange
	var err error
	var entries *webdb.RainEntriesMm
	var resp []byte

	if dates, err = getToFromTotal(res); err != nil {
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
