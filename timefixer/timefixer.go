// Package timefixer fixes www.energiatili.fi timestamps.
//
// The timestamps are mangled as follows (which this fixer reverses):
//   1. Take a timestamp in Finnish local time (EET/EEST)
//   2. Ignore time zone information from local time and interpret time as UTC
//   3. Represent time in Unix Epoch style, milliseconds from 1970-01-01 00:00:00 UTC
// The reversing must keep state since on DST change either one hour is missing or repeated.
package timefixer

import (
	"log"
	"time"
)

// TimeFixer stores information about the previous timestamp needed
// for correction
type TimeFixer struct {
	prev       time.Time
	adjustment time.Duration
}

// ParseBrokenTime parses timestamps statefully and fixes timezone
func (t *TimeFixer) ParseBrokenTime(localTs float64) (ts time.Time, err error) {
	// Fix time
	realTs := time.Unix(int64(localTs/1000), 0).Add(-2 * time.Hour).UTC()

	diff := realTs.Sub(t.prev)
	t.prev = realTs
	switch {
	case t.prev.IsZero() == true:
		//
	case diff == 0*time.Hour:
		// DST backward; following data in standard time
		log.Printf("WARN: Detected DST backward shift %v\n", realTs)
		t.adjustment = 0 * time.Hour
		//realTs = realTs.Add(-1 * time.Hour)
	case diff == 1*time.Hour:
		// Normal case
	case diff == 2*time.Hour:
		// DST forward, following data in summer time
		t.adjustment = 1 * time.Hour
		log.Printf("WARN: Detected DST forward shift %v\n", realTs)
	case diff == 10*time.Hour || diff == 16*time.Hour:
		// Day meter vs. Night meter gap
	default:
		_, offset := realTs.Local().Zone()
		t.adjustment = time.Duration(offset/3600-2) * time.Hour
		log.Printf("WARN: Records missing, now reading %v (got diff=%v)\n", realTs, diff)
	}

	realTs = realTs.Add(-t.adjustment)

	return realTs, nil
}
