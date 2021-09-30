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
	args, err := ParseQuery(raw)
	if err != nil {
		handler.badRequest(w, err)
	}

	since := handler.statusDurationDefault
	_, ok := args["since"]
	if ok {
		var asNum int
		if asNum, err = strconv.Atoi(args["since"].(string)); err != nil {
			handler.badRequest(w, err)
			return
		}
		since = time.Second * time.Duration(asNum)
	}
	var isUp bool
	if isUp, err = dbQuery(since); err != nil {
		handler.internalServiceError(w, err)
		return
	}
	resp, err := json.Marshal(map[string]interface{}{responseKey: isUp})
	if err != nil {
		logrus.Error(err)
	}
	handler.genericJSONHandler(resp, w, res)
}
