package templates

import (
	"time"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
)

// WeatherData is the raw data sent to the templates
type WeatherData struct {
	// standard rain measurements
	HourRainIn           string
	TwoHourRainIn        string
	SixHourRainIn        string
	TwentyFourHourRainIn string
	SevenDayRainIn       string
	ThirtyDayRainIn      string
	YearlyRainIn         string

	// metric rain measurements
	HourRainMm           string
	TwoHourRainMm        string
	SixHourRainMm        string
	TwentyFourHourRainMm string
	SevenDayRainMm       string
	ThirtyDayRainMm      string
	YearlyRainMm         string

	// various database lookups
	TempF         int
	TempC         int
	LastRain      string
	GatewayStatus string
	SensorStatus  string

	// simple datetime
	LastUpdate string
	Year       int

	// consts
	HourIndicator       string
	DayIndicator        string
	InchIndicator       string
	MillimeterIndicator string
	YearIndicator       string
}

const ErrorFloatString = "-999.9"
const ErrorInt = -999
const ErrorTimestamp = "ERROR getting timestamp"
const ErrorStatus = "ERROR getting status"

var BaseWeatherData = WeatherData{
	HourRainIn:           ErrorFloatString,
	TwoHourRainIn:        ErrorFloatString,
	SixHourRainIn:        ErrorFloatString,
	TwentyFourHourRainIn: ErrorFloatString,
	SevenDayRainIn:       ErrorFloatString,
	ThirtyDayRainIn:      ErrorFloatString,
	YearlyRainIn:         ErrorFloatString,

	HourRainMm:           ErrorFloatString,
	TwoHourRainMm:        ErrorFloatString,
	SixHourRainMm:        ErrorFloatString,
	TwentyFourHourRainMm: ErrorFloatString,
	SevenDayRainMm:       ErrorFloatString,
	ThirtyDayRainMm:      ErrorFloatString,
	YearlyRainMm:         ErrorFloatString,

	TempF:         ErrorInt,
	TempC:         ErrorInt,
	LastRain:      ErrorTimestamp,
	GatewayStatus: ErrorStatus,
	SensorStatus:  ErrorStatus,

	LastUpdate: time.Now().Format(configkey.PrettyTimeFormat),
	Year:       time.Now().Year(),

	HourIndicator:       "hr:",
	DayIndicator:        "d:",
	YearIndicator:       "ytd:",
	InchIndicator:       "in",
	MillimeterIndicator: "mm",
}
