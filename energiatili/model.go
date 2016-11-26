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
				Data           []DataPoint
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

// MarshalJSON returns d as the JSON encoding of d.
func (d *DataPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal([]float64{float64(d.Time.Unix() * 1000), d.Kwh})
}

// UnmarshalJSON sets d to a copy of data.
func (d *DataPoint) UnmarshalJSON(data []byte) error {
	var decoded [2]float64
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	d.Time = parseEnergiatiliTime(decoded[0])
	d.Kwh = decoded[1]
	return nil
}

// parseEnergiatiliTime decodes "unixMillis" ignoring time zone and cast to Helsinki time
func parseEnergiatiliTime(t float64) time.Time {
	ts := time.Unix(int64(t/1000), 0).UTC()
	year, _, day := ts.Date()
	month := ts.Month()
	hour, min, sec := ts.Clock()
	return time.Date(year, month, day, hour, min, sec, 0, helsinki)
}

var helsinki *time.Location

func init() {
	var err error
	helsinki, err = time.LoadLocation("Europe/Helsinki")
	if err != nil {
		panic(err)
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

// DataPoints returns all the consumption readings in the report.
// The records are valid even if ErrorMissingRecord is returned to indicate
// gaps in data.
func (c *ConsumptionReport) DataPoints() (points []DataPoint, err error) {
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

			// kwh := raw[1]
			// points = append(points, DataPoint{Time: ts, Kwh: kwh})
		}
	}
	if missingRecords {
		return points, ErrorMissingRecord
	}
	return points, nil
}
