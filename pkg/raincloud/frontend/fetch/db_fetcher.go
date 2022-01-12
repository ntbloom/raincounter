package fetch

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ntbloom/raincounter/pkg/config/configkey"

	"github.com/ntbloom/raincounter/pkg/raincloud/frontend/templates"
	"github.com/sirupsen/logrus"

	"github.com/ntbloom/raincounter/pkg/raincloud/webdb"
)

type DataFetcher struct {
	query webdb.DBQuery
	data  templates.WeatherData
	sync.Mutex
}

func NewDataFetcher() *DataFetcher {
	return &DataFetcher{
		query: webdb.NewPGConnector(),
		data:  templates.BaseWeatherData,
		Mutex: sync.Mutex{},
	}
}

func (d *DataFetcher) Fetch() interface{} {
	now := time.Now()
	var wg sync.WaitGroup

	// time specific functions
	type callableWithTimestamp func(time.Time)
	for _, v := range []callableWithTimestamp{
		d.getHourRain,
		d.getTwoHourRain,
		d.getSixHourRain,
		d.getTwentyFourHourRain,
		d.getSevenDayRain,
		d.getThirtyDayRain,
		d.getYearTotalRain,
	} {
		wg.Add(1)
		get := v
		go func() {
			d.Lock()
			defer d.Unlock()
			defer wg.Done()
			get(now)
		}()
	}

	// parameterless functions
	type callable func()
	for _, v := range []callable{
		d.getCurrentTemp,
		d.getLastRain,
		d.getSensorStatus,
		d.getGatewayStatus,
	} {
		wg.Add(1)
		get := v
		go func() {
			d.Lock()
			defer d.Unlock()
			defer wg.Done()
			get()
		}()
	}
	wg.Wait()
	return d.data
}

func formatFloatFromDatabase(val float64, err error) (std string, metric string) {
	if err != nil {
		logrus.Error(err)
		return templates.ErrorFloatString, templates.ErrorFloatString
	}
	format := "%.2f"
	metric = fmt.Sprintf(format, val)
	const toInches = 0.0393701
	std = fmt.Sprintf(format, val*toInches)
	return
}

func (d *DataFetcher) getRainSince(now time.Time, duration time.Duration) (std string, metric string) {
	val, err := d.query.TotalRainMMFrom(now.Add(duration), now)
	return formatFloatFromDatabase(val, err)
}

func (d *DataFetcher) getHourRain(now time.Time) {
	std, met := d.getRainSince(now, time.Hour*-1)
	d.data.HourRainIn = std
	d.data.HourRainMm = met
}

func (d *DataFetcher) getTwoHourRain(now time.Time) {
	std, met := d.getRainSince(now, time.Hour*-2)
	d.data.TwoHourRainIn = std
	d.data.TwoHourRainMm = met
}

func (d *DataFetcher) getSixHourRain(now time.Time) {
	std, met := d.getRainSince(now, time.Hour*-6)
	d.data.SixHourRainIn = std
	d.data.SixHourRainMm = met
}

func (d *DataFetcher) getTwentyFourHourRain(now time.Time) {
	std, met := d.getRainSince(now, time.Hour*-24)
	d.data.TwentyFourHourRainIn = std
	d.data.TwentyFourHourRainMm = met
}

func (d *DataFetcher) getSevenDayRain(now time.Time) {
	const seven = 7
	std, met := d.getRainSince(now, time.Hour*-24*seven)
	d.data.SevenDayRainIn = std
	d.data.SevenDayRainMm = met
}

func (d *DataFetcher) getThirtyDayRain(now time.Time) {
	const thirty = 30
	std, met := d.getRainSince(now, time.Hour*-24*thirty)
	d.data.ThirtyDayRainIn = std
	d.data.ThirtyDayRainMm = met
}

func (d *DataFetcher) getYearTotalRain(now time.Time) {
	val, err := d.query.TotalRainMMFrom(time.Date(now.Year(), 0, 0, 0, 0, 0, 0, time.UTC), now)
	std, met := formatFloatFromDatabase(val, err)
	d.data.YearlyRainIn = std
	d.data.YearlyRainMm = met
}

func (d *DataFetcher) getCurrentTemp() {
	tempC, err := d.query.GetLastTempC()
	if err != nil {
		logrus.Errorf("error getting current temp: %s", err)
		return
	}
	const thirtytwo = 32
	tempF := int(math.Round(float64(tempC)*9/5) + thirtytwo)
	d.data.TempC = tempC
	d.data.TempF = tempF
}

func (d *DataFetcher) getLastRain() {
	date, err := d.query.GetLastRainTime()
	if err != nil {
		logrus.Errorf("error getting last rain: %s", err)
		return
	}
	d.data.LastRain = date.Format(configkey.PrettyTimeFormat)
}

type callable func(time.Duration) (bool, error)

func (d *DataFetcher) getStatus(c callable) (string, error) {
	since := time.Since(time.Now().Add(time.Minute * -5))
	up, err := c(since)
	if err != nil {
		return "", err
	}
	if up {
		return "up", nil
	} else {
		return "down", nil
	}
}

func (d *DataFetcher) getGatewayStatus() {
	val, err := d.getStatus(d.query.IsGatewayUp)
	if err != nil {
		logrus.Errorf("error querying gateway status: %s", err)
		return
	}
	d.data.GatewayStatus = val
}

func (d *DataFetcher) getSensorStatus() {
	val, err := d.getStatus(d.query.IsSensorUp)
	if err != nil {
		logrus.Errorf("error querying sensor status: %s", err)
		return
	}
	d.data.SensorStatus = val

}
