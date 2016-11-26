package energiatili

import (
	"encoding/json"
	"io"
)

// ConsumptionReport is the structure in 'var model' of Energiatili
type ConsumptionReport struct {
	Hours struct {
		Temperature struct {
			Data []Record
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
				Data           []Record
			}
		}
	}
}

// FromJSON parses consumption report JSON
func FromJSON(jsonData io.Reader) (c *ConsumptionReport, err error) {
	decoder := json.NewDecoder(jsonData)
	if err := decoder.Decode(&c); err != nil {
		return nil, err
	}
	return c, err
}

// Records returns all the consumption readings in the report.
// The records are valid even if ErrorMissingRecord is returned to indicate
// gaps in data.
func (c *ConsumptionReport) Records() (points []Record, err error) {
	// fixer := TimeFixer{}
	missingRecords := false
	for _, cons := range c.Hours.Consumptions {
		for _, p := range cons.Series.Data {
			points = append(points, p)
			// ts, err := fixer.ParseBrokenTime(raw[0])
			// if err != nil {
			// 	if err != ErrorMissingRecord {
			// 		return nil, err
			// 	}
			// 	missingRecords = true
			// }

			// value := raw[1]
			// points = append(points, Record{Time: ts, Value: value})
		}
	}
	if missingRecords {
		return points, ErrorMissingRecord
	}
	return points, nil
}
