package energiatili

import (
	"encoding/json"
	"fmt"
	"time"
)

// Record is a timestamped value (consumption or temperature)
type Record struct {
	Timestamp time.Time
	Value     float64
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

// MarshalJSON returns d as the JSON encoding of d.
func (r Record) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%d,%f]", formatEnergiatiliTime(r.Timestamp), r.Value)), nil
}

// UnmarshalJSON sets d to a copy of data.
func (r *Record) UnmarshalJSON(data []byte) error {
	var decoded [2]float64
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	r.Timestamp = parseEnergiatiliTime(decoded[0])
	r.Value = decoded[1]
	return nil
}

// ByTime is used to sort records by their Time field.
type ByTime []Record

func (b ByTime) Len() int           { return len(b) }
func (b ByTime) Less(i, j int) bool { return b[i].Timestamp.Before(b[j].Timestamp) }
func (b ByTime) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

// Records implements notz.Interface for notz.FixDST.
type Records []Record

func (r Records) Len() int                     { return len(r) }
func (r Records) Time(i int) time.Time         { return r[i].Timestamp }
func (r Records) SetTime(i int, new time.Time) { r[i].Timestamp = new }

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
