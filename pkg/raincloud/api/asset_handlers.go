package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

/* HANDLE THE ASSETS (SENSOR AND GATEWAY) AND EVENT MESSAGES */

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
	args := ParseQuery(raw)

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
