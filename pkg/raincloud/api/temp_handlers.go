package api

import (
	"encoding/json"
	"net/http"

	"github.com/ntbloom/raincounter/pkg/raincloud/webdb"
	"github.com/sirupsen/logrus"
)

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

// handle request for temperature data
func (handler restHandler) handleTemp(w http.ResponseWriter, res *http.Request) {
	var dates *dateRange
	var err error
	var entries *webdb.TempEntriesC
	var resp []byte

	if dates, err = getToFromTotal(res); err != nil {
		handler.badRequest(w, err)
		return
	}
	badRequest := dates == nil || !dates.fromOk
	if badRequest {
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
