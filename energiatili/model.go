package energiatili

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
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

// Record is a timestamped value (consumption or temperature)
type Record struct {
	Time  time.Time
	Value float64
}

// MarshalJSON returns d as the JSON encoding of d.
func (d Record) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%d,%f]", formatEnergiatiliTime(d.Time), d.Value)), nil
}

// UnmarshalJSON sets d to a copy of data.
func (d *Record) UnmarshalJSON(data []byte) error {
	var decoded [2]float64
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	d.Time = parseEnergiatiliTime(decoded[0])
	d.Value = decoded[1]
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

// formatEnergiatiliTime encodes Helsinki time "unixMillis" ignoring time zone
func formatEnergiatiliTime(t time.Time) int64 {
	tHelsinki := t.In(helsinki)
	year, _, day := tHelsinki.Date()
	month := tHelsinki.Month()
	hour, min, sec := tHelsinki.Clock()
	return time.Date(year, month, day, hour, min, sec, 0, utc).Unix() * 1000
}

var utc, helsinki *time.Location

func init() {
	var err error
	helsinki, err = time.LoadLocation("Europe/Helsinki")
	if err != nil {
		panic(err)
	}
	utc, err = time.LoadLocation("UTC")
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
