package fetcher

import (
	"errors"
	"time"

	"github.com/joneskoo/etget/timefixer"
)

// ConsumptionReport is the structure in 'var model' of Energiatili
type ConsumptionReport struct {
	Hours struct {
		Temperature struct {
			RawData [][2]float64 `json:"Data"`
		}
		Step struct {
			TimeZoneInfo struct {
				ID              string `json:"Id"`
				DisplayName     string
				StandardName    string
				DaylightName    string
				BaseUtcOffset   string
				AdjustmentRules interface{}
			}
		}
		Consumptions []struct {
			TariffTimeZoneID          int `json:"TariffTimeZoneId"`
			TariffTimeZoneName        string
			TariffTimeZoneDescription string
			Series                    struct {
				ReadingCounter int
				Name           string
				Resolution     string
				RawData        [][2]float64 `json:"Data"`
				Data           []DataPoint  `json:"-"`
			}
		}
	}
}

// DataPoint is the parsed format of a single record, result
// of running Update()
type DataPoint struct {
	Time time.Time
	Kwh  float64
}

// Update processes input to a easier-to-parse format.
// Namely RawData under Consumptions is parsed to Data with UTC timestamps.
func (c *ConsumptionReport) Update() (err error) {
	fixer := timefixer.TimeFixer{}
	for i, cons := range c.Hours.Consumptions {
		count := len(cons.Series.RawData)
		c.Hours.Consumptions[i].Series.Data = make([]DataPoint, count)
		for j, raw := range cons.Series.RawData {
			if len(raw) != 2 {
				return errors.New("Invalid data") //FIXME
			}
			c.Hours.Consumptions[i].Series.Data[j].Time, err = fixer.ParseBrokenTime(raw[0])
			c.Hours.Consumptions[i].Series.Data[j].Kwh = raw[1]
		}
	}
	return nil
}
