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
	return points, nil
}
