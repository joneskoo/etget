// Package timefixer fixes www.energiatili.fi timestamps.
//
// The timestamps are mangled as follows (which this fixer reverses):
//   1. Take a timestamp in Finnish local time (EET/EEST)
//   2. Ignore time zone information from local time and interpret time as UTC
//   3. Represent time in Unix Epoch style, milliseconds from 1970-01-01 00:00:00 UTC
// The reversing must keep state since on DST change either one hour is missing or repeated.
package timefixer

import (
	"fmt"
	"log"
	"time"
)

// TimeFixer stores information about the previous timestamp needed
// for correction
type TimeFixer struct {
	prev time.Time
}

// ParseBrokenTime parses timestamps statefully and fixes timezone
func (t *TimeFixer) ParseBrokenTime(localTs float64) (ts time.Time, err error) {
	// Decode "unix" time and ignore time zone
	realTs := time.Unix(int64(localTs/1000), 0).UTC()
	year, _, day := realTs.Date()
	month := realTs.Month()
	hour, min, sec := realTs.Clock()
	helsinki, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		return time.Time{}, err
	}

	// Interpret clock time as local time
	realTs = time.Date(year, month, day, hour, min, sec, 0, helsinki)

	// Compare to previous record; adjust for DST gap
	diff := realTs.Sub(t.prev)
	switch diff {
	case 0 * time.Hour:
		// DST backward; following data in standard time
		// This should not happen because local time can be decoded uniquely
		// to UTC as local time jumps ahead an extra hour
		realTs = realTs.Add(1 * time.Hour)
		fmt.Printf("WARN: Compensating DST backward shift %v, this should not happen\n", realTs)
	case 1 * time.Hour:
		// Normal case
	case 2 * time.Hour:
		// DST forward, following data in summer time
		realTs = realTs.Add(-1 * time.Hour)
		fmt.Printf("WARN: Compensating DST forward shift %v\n", realTs)
	default:
		if !t.prev.IsZero() {
			log.Printf("WARN: Records missing, now reading %v (got diff=%v)\n", realTs, diff)
		}
	}
	t.prev = realTs
	return realTs, nil
}
