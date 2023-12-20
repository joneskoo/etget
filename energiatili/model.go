package energiatili

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/joneskoo/etget/notz"
)

// ConsumptionReport is the structure in 'var model' of Energiatili
type ConsumptionReport struct {
	IsValid                bool   `json:"IsValid"`
	HasTemperatureSeries   bool   `json:"HasTemperatureSeries"`
	HasReactivePowerSeries bool   `json:"HasReactivePowerSeries"`
	PowerUnit              string `json:"PowerUnit"`
	DataInterval           struct {
		Duration   string  `json:"Duration"`
		Start      Date    `json:"Start"`
		Stop       Date    `json:"Stop"`
		StartValue Date    `json:"StartValue"`
		StopValue  Date    `json:"StopValue"`
		TotalYears float64 `json:"TotalYears"`
		TotalDays  float64 `json:"TotalDays"`
		TotalHours float64 `json:"TotalHours"`
	} `json:"DataInterval"`
	Hours            Report           `json:"Hours"`
	Days             Report           `json:"Days"`
	Weeks            Report           `json:"Weeks"`
	Months           Report           `json:"Months"`
	Years            Report           `json:"Years"`
	NetworkPriceList NetworkPriceList `json:"NetworkPriceList"`
	EnergyTaxes      struct {
		Code  string `json:"Code"`
		Name  string `json:"Name"`
		Taxes []struct {
			TotalTaxWithVAT float64 `json:"TotalTaxWithVAT"`
			TotalTaxNoVAT   float64 `json:"TotalTaxNoVAT"`
			StartDate       Date    `json:"StartDate"`
			EndDate         Date    `json:"EndDate"`
			StartDateTicks  int64   `json:"StartDateTicks"`
			EndDateTicks    int64   `json:"EndDateTicks"`
		} `json:"Taxes"`
	} `json:"EnergyTaxes"`
	Vat []struct {
		Tax   float64 `json:"Tax"`
		Start Date    `json:"Start"`
		Stop  Date    `json:"Stop"`
	} `json:"Vat"`
	UseReadings                                   bool    `json:"UseReadings"`
	HourlyDataAvailable                           bool    `json:"HourlyDataAvailable"`
	UtilityType                                   int     `json:"UtilityType"`
	CustomerType                                  int     `json:"CustomerType"`
	LowTemperature                                float64 `json:"LowTemperature"`
	HighTemperature                               float64 `json:"HighTemperature"`
	HasProviderHourlyValuesSeparateLoadingSupport bool    `json:"HasProviderHourlyValuesSeparateLoadingSupport"`
	HasHourlyValues                               bool    `json:"HasHourlyValues"`
	FullConsumptionLoaded                         bool    `json:"FullConsumptionLoaded"`
	ContractTypeID                                int     `json:"ContractTypeId"`
	HasWaterHourlyValues                          bool    `json:"HasWaterHourlyValues"`
}

func (c *ConsumptionReport) Records() (points []Record, err error) {
	for _, cons := range c.Hours.Consumptions {
		for _, p := range cons.Series.Data {
			points = append(points, p)
		}
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Timestamp.Before(points[j].Timestamp)
	})
	notz.FixDST(records(points))
	points = trimTrailingZeros(points)
	return points, nil
}

type records []Record

func (r records) Len() int                     { return len(r) }
func (r records) Time(i int) time.Time         { return r[i].Timestamp }
func (r records) SetTime(i int, new time.Time) { r[i].Timestamp = new }

type Date time.Time

func (d *Date) UnmarshalJSON(data []byte) error {
	var tmp float64
	data = bytes.TrimPrefix(data, []byte(`"/Date(`))
	data = bytes.TrimSuffix(data, []byte(`)/"`))
	if err := json.Unmarshal(data, &tmp); err != nil {
		fmt.Printf("failed to parse: %q\n", data)
		return err
	}
	*d = Date(parseEnergiatiliTime(tmp))
	return nil
}

type Duration string // "02:00:00"

type Report struct {
	PMax                PMax          `json:"PMax"`
	Step                Step          `json:"Step"`
	Consumptions        []Consumption `json:"Consumptions"`
	Productions         []Consumption `json:"Productions"`
	SalesConsumptions   []Consumption `json:"SalesConsumptions"`
	Temperature         Series        `json:"Temperature"`
	ActivePower         Series        `json:"ActivePower"`
	ConsumptionStatuses Series        `json:"ConsumptionStatuses"`
	MeterReadings       Series        `json:"MeterReadings"`
}

type TimeZoneInfo struct {
	ID                         string           `json:"Id"`
	DisplayName                string           `json:"DisplayName"`
	StandardName               string           `json:"StandardName"`
	DaylightName               string           `json:"DaylightName"`
	BaseUtcOffset              string           `json:"BaseUtcOffset"`
	AdjustmentRules            []AdjustmentRule `json:"AdjustmentRules"`
	SupportsDaylightSavingTime bool             `json:"SupportsDaylightSavingTime"`
}

type AdjustmentRule struct {
	DateStart               Date     `json:"DateStart"`
	DateEnd                 Date     `json:"DateEnd"`
	DaylightDelta           Duration `json:"DaylightDelta"`
	DaylightTransitionStart struct {
		TimeOfDay       Date `json:"TimeOfDay"`
		Month           int  `json:"Month"`
		Week            int  `json:"Week"`
		Day             int  `json:"Day"`
		DayOfWeek       int  `json:"DayOfWeek"`
		IsFixedDateRule bool `json:"IsFixedDateRule"`
	} `json:"DaylightTransitionStart"`
	DaylightTransitionEnd struct {
		TimeOfDay       Date `json:"TimeOfDay"`
		Month           int  `json:"Month"`
		Week            int  `json:"Week"`
		Day             int  `json:"Day"`
		DayOfWeek       int  `json:"DayOfWeek"`
		IsFixedDateRule bool `json:"IsFixedDateRule"`
	} `json:"DaylightTransitionEnd"`
	BaseUtcOffsetDelta Duration `json:"BaseUtcOffsetDelta"`
}

type PMax struct {
	Item1 Date    `json:"Item1"`
	Item2 float64 `json:"Item2"`
}

type NetworkPriceList struct {
	PriceListName        string `json:"PriceListName"`
	ActiveTarificationID int    `json:"ActiveTarificationId"`
	Tarifications        []struct {
		TarificationID int  `json:"TarificationId"`
		StartTime      Date `json:"StartTime"`
		EndTime        Date `json:"EndTime"`
	} `json:"Tarifications"`
	TimeBasedEnergyDayPrices []struct {
		ProductComponentTypeCode string  `json:"ProductComponentTypeCode"`
		StartTime                Date    `json:"StartTime"`
		EndTime                  Date    `json:"EndTime"`
		PriceWithVat             float64 `json:"PriceWithVat"`
		PriceNoVat               float64 `json:"PriceNoVat"`
	} `json:"TimeBasedEnergyDayPrices"`
	TimeBasedEnergyNightPrices []struct {
		ProductComponentTypeCode string  `json:"ProductComponentTypeCode"`
		StartTime                Date    `json:"StartTime"`
		EndTime                  Date    `json:"EndTime"`
		PriceWithVat             float64 `json:"PriceWithVat"`
		PriceNoVat               float64 `json:"PriceNoVat"`
	} `json:"TimeBasedEnergyNightPrices"`
	IsSpotPriceList bool `json:"IsSpotPriceList"`
}

type Step struct {
	Type         int          `json:"Type"`
	TimeZoneInfo TimeZoneInfo `json:"TimeZoneInfo"`
	StepLength   string       `json:"StepLength,omitempty"`
	Start        Date         `json:"Start"`
	Stop         Date         `json:"Stop"`
	StepCount    int          `json:"StepCount"`
}

type Consumption struct {
	TariffTimeZoneID          int    `json:"TariffTimeZoneId"`
	TariffTimeZoneName        string `json:"TariffTimeZoneName"`
	TariffTimeZoneDescription string `json:"TariffTimeZoneDescription"`
	Series                    Series `json:"Series"`
}

type Series struct {
	ReadingCounter int    `json:"ReadingCounter"`
	Name           string `json:"Name"`
	Resolution     string `json:"Resolution"`
	Data           Data   `json:"Data"`
	Start          Date   `json:"Start"`
	Stop           Date   `json:"Stop"`
	Step           Step   `json:"Step"`
	DataCount      int    `json:"DataCount"`
	Type           string `json:"Type"`
	Unit           string `json:"Unit"`
}

type Data []Record

// Record is a timestamped value (consumption or temperature)
type Record struct {
	Timestamp time.Time
	Value     float64
}

// UnmarshalJSON sets r to a copy of data.
func (r *Record) UnmarshalJSON(data []byte) error {
	var decoded [2]float64
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	r.Timestamp = parseEnergiatiliTime(decoded[0])
	r.Value = decoded[1]
	return nil
}

func trimTrailingZeros(p []Record) []Record {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i].Value != 0 {
			return p[0 : i+1]
		}
	}
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
