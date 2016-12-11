package energiatili

import (
	"sort"

	"github.com/joneskoo/etget/notz"
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

// Records returns all the consumption readings in the report.
func (c ConsumptionReport) Records() (points []Record, err error) {
	for _, cons := range c.Hours.Consumptions {
		for _, p := range cons.Series.Data {
			points = append(points, p)
		}
	}
	sort.Sort(byTime(points))
	notz.FixDST(records(points))
	points = trimTrailingZeros(points)
	return points, nil
}

func trimTrailingZeros(p []Record) []Record {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i].Value != 0 {
			return p[0 : i+1]
		}
	}
	return nil
}
