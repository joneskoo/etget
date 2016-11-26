package energiatili

import (
	"encoding/json"
	"fmt"
	"time"
)

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

// Record is a timestamped value (consumption or temperature)
type Record struct {
	Timestamp time.Time
	Value     float64
}

// MarshalJSON returns r as the JSON encoding of r.
func (r Record) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%d,%f]", formatEnergiatiliTime(r.Timestamp), r.Value)), nil
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

// byTime is used to sort records by their Time field.
type byTime []Record

func (b byTime) Len() int           { return len(b) }
func (b byTime) Less(i, j int) bool { return b[i].Timestamp.Before(b[j].Timestamp) }
func (b byTime) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

// records implements notz.Interface for notz.FixDST.
type records []Record

func (r records) Len() int                     { return len(r) }
func (r records) Time(i int) time.Time         { return r[i].Timestamp }
func (r records) SetTime(i int, new time.Time) { r[i].Timestamp = new }

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
