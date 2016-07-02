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
	"strconv"
	"time"
)

// TimeFixer stores information about the previous timestamp needed
// for correction
type TimeFixer struct {
	prev       time.Time
	adjustment time.Duration
}

// ParseBrokenTime parses timestamps statefully and fixes timezone
func (t *TimeFixer) ParseBrokenTime(s string) (ts time.Time, err error) {
	localTs, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	// Fix time
	realTs := time.Unix(localTs, 0).Add(-2 * time.Hour).UTC()

	diff := realTs.Sub(t.prev)
	t.prev = realTs
	switch diff {
	case 0 * time.Hour:
		// DST backward; following data in standard time
		log.Printf("WARN: Detected DST backward shift %v\n", realTs)
		t.adjustment = 0 * time.Hour
		//realTs = realTs.Add(-1 * time.Hour)
	case 1 * time.Hour:
		// Normal case
	case 2 * time.Hour:
		// DST forward, following data in summer time
		t.adjustment = 1 * time.Hour
		log.Printf("WARN: Detected DST forward shift %v\n", realTs)
	default:
		_, offset := realTs.Local().Zone()
		t.adjustment = time.Duration(offset/3600-2) * time.Hour
		log.Printf("WARN: Records missing, now reading %v\n", realTs)
	}

	realTs = realTs.Add(-t.adjustment)

	return realTs, nil
}
