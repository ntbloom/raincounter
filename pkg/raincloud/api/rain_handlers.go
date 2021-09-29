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
	var resp []byte

	if dates, err = getToFromTotal(res); err != nil {
		handler.badRequest(w, err)
		return
	}
	if dates.total {
		// only show the total rain from the period
		var amt float64
		if dates.toOk {
			if amt, err = handler.db.TotalRainMMFrom(dates.from, dates.to); err != nil {
				handler.internalServiceError(w, err)
				return
			}
		} else {
			if amt, err = handler.db.TotalRainMMSince(dates.from); err != nil {
				handler.internalServiceError(w, err)
				return
			}
		}
		if resp, err = json.Marshal(map[string]float64{"amount": amt}); err != nil {
			handler.internalServiceError(w, err)
			return
		}
		logrus.Debug(amt)
	} else {
		// give a struct with the amounts and timestamps of each rain event
		var entries *webdb.RainEntriesMm
		if dates.toOk {
			if entries, err = handler.db.GetRainMMFrom(dates.from, dates.to); err != nil {
				handler.internalServiceError(w, err)
				return
			}
		} else {
			if entries, err = handler.db.GetRainMMSince(dates.from); err != nil {
				handler.internalServiceError(w, err)
				return
			}
		}
		if resp, err = json.Marshal(entries); err != nil {
			handler.internalServiceError(w, err)
			return
		}
	}
	handler.genericJSONHandler(resp, w, res)
}
