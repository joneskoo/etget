package fetcher

import (
	"encoding/json"
	"time"
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

// FromJSON parses consumption report JSON
func FromJSON(jsonData []byte) (c ConsumptionReport, err error) {
	if err = json.Unmarshal(jsonData, &c); err != nil {
		return
	}
	err = c.update()
	return
}

// update parses raw timestamps to native time.Time
func (c *ConsumptionReport) update() (err error) {
	fixer := TimeFixer{}
	for i, cons := range c.Hours.Consumptions {
		count := len(cons.Series.RawData)
		c.Hours.Consumptions[i].Series.Data = make([]DataPoint, count)
		for j, raw := range cons.Series.RawData {
			c.Hours.Consumptions[i].Series.Data[j].Time, err = fixer.ParseBrokenTime(raw[0])
			c.Hours.Consumptions[i].Series.Data[j].Kwh = raw[1]
		}
	}
	return nil
}
