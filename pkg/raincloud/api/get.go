package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ntbloom/raincounter/pkg/config/configkey"

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
	handler.writeJSONResponse(resp, w, res)
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
	badRequest := dates == nil || !dates.fromOk
	if badRequest {
		handler.badRequest(w, fmt.Errorf("no arguments provided"))
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
	handler.writeJSONResponse(resp, w, res)
}

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
	handler.writeJSONResponse(resp, w, res)
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
	handler.writeJSONResponse(resp, w, res)
}

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
	args, err := ParseRestQuery(raw)
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
	handler.writeJSONResponse(resp, w, res)
}

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
	handler.writeJSONResponse(payload, w, res)
}

/* HELPER METHODS FOR GET REQUESTS */

// handles generic JSON messages. fails if the request does not specify application/json
func (handler restHandler) writeJSONResponse(payload []byte, w http.ResponseWriter, res *http.Request) {
	// allow CORS responses
	w.Header().Set("Access-Control-Allow-Origin", "*")

	encoding := res.Header.Get(contentType)
	if encoding != appJSON {
		handler.unsupportedMedia(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		logrus.Error(err)
	}
	logrus.Infof("headers sent=%s", w.Header())
}

// ParseRestQuery breaks the restful part of the API into a map
func ParseRestQuery(raw string) (map[string]interface{}, error) {
	var err error
	if raw == "" {
		err = fmt.Errorf("empty arguments")
		return nil, err
	}
	result := make(map[string]interface{})
	args := strings.Split(raw, "&")
	for _, arg := range args {
		keys := strings.Split(arg, "=")
		if len(keys) < 2 {
			err = fmt.Errorf("illegal REST argument: %s", arg)
			return nil, err
		}
		result[keys[0]] = keys[1]
	}
	return result, nil
}

// dateRange is a parsed struct of JSON data with to and from timestamp
type dateRange struct {
	toOk    bool
	fromOk  bool
	totalOk bool
	to      time.Time
	from    time.Time
	total   bool
}

// get args from the rest API
func getToFromTotal(res *http.Request) (*dateRange, error) {
	var err error
	var args map[string]interface{}

	if args, err = ParseRestQuery(res.URL.RawQuery); err != nil {
		return nil, fmt.Errorf("unparseable arguments")
	}

	_, fromOk := args["from"]
	if !fromOk {
		return nil, err
	}
	var from time.Time
	if from, err = parseTimestamp(args["from"].(string)); err != nil {
		return nil, err
	}

	var to time.Time
	_, toOk := args["to"]
	if toOk {
		if to, err = parseTimestamp(args["to"].(string)); err != nil {
			return nil, err
		}
	}
	total := false
	_, totalOk := args["total"]
	if totalOk {
		if total, err = strconv.ParseBool(args["total"].(string)); err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return &dateRange{
		toOk:    toOk,
		to:      to,
		fromOk:  fromOk,
		from:    from,
		totalOk: totalOk,
		total:   total,
	}, nil
}

// parseTimestamp gets a timestamp from the raw interface
func parseTimestamp(raw string) (time.Time, error) {
	var stamp time.Time
	var err error

	attempt := func(format string) time.Time {
		var parsed time.Time
		if parsed, err = time.Parse(format, raw); err != nil {
			logrus.Errorf("doesn't match %s: %s", format, err)
			return time.Time{}
		}
		return parsed
	}
	stamp = attempt(configkey.TimestampFormat)

	return stamp, err
}
