package energiatili

import (
	"encoding/json"
	"io"
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
				Data           [][2]float64
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
func FromJSON(jsonData io.Reader) (c *ConsumptionReport, err error) {
	decoder := json.NewDecoder(jsonData)
	if err := decoder.Decode(&c); err != nil {
		return nil, err
	}
	return c, err
}

// DataPoints returns all the consumption readings in the report.
// The records are valid even if ErrorMissingRecord is returned to indicate
// gaps in data.
func (c *ConsumptionReport) DataPoints() (points []DataPoint, err error) {
	fixer := TimeFixer{}
	missingRecords := false
	for _, cons := range c.Hours.Consumptions {
		for _, raw := range cons.Series.Data {
			ts, err := fixer.ParseBrokenTime(raw[0])
			if err != nil {
				if err != ErrorMissingRecord {
					return nil, err
				}
				missingRecords = true
			}

			kwh := raw[1]
			points = append(points, DataPoint{Time: ts, Kwh: kwh})
		}
	}
	if missingRecords {
		return points, ErrorMissingRecord
	}
	return points, nil
}
