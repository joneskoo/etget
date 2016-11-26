// Package notz is a post-processing fix for DST-confused timestamps.
// Namely, in the case where there are hourly record timestamps but the
// timestamp does not include the timezone (standard or DST), this package
// can restore the corrupted timestamps.
//
// Where possible, it is obviously preferrable to use UTC timestamps or
// if local time must be used, include the timezone name/offset with the time.
// This package is only intended to recover data from external sources.
//
// For example a transition from summer time to winter time (Helsinki):
//
//     Sun Oct 25 02:00:00 EEST 2015
//     Sun Oct 25 03:00:00 EEST 2015
//     Sun Oct 25 03:00:00 EET 2015
//     Sun Oct 25 04:00:00 EET 2015
//
// Note that the local time is identical, with difference only in
// time zone name. If we parse "Sun Oct 25 03:00:00 2015",
// we get "Sun Oct 25 03:00:00 EET 2015".
// If we get two identical timestamps, we assume that the first
// was mis-interpreted and restore it by subtracting 1 hour.
package notz

import "time"

// TimeSetter interface must be implemented by values used with FixDST.
type TimeSetter interface {
	// Time retrieves the timestamp value to be fixed
	Time() time.Time

	// SetTime sets the timestamp value to be fixed
	SetTime(time.Time)
}

// FixDST fixes DST ambiguoity in a slice of values.
// All values must be hours with no minutes, seconds or nanoseconds part
// in sequential order.
func FixDST(times []TimeSetter) {
	var prev TimeSetter
	for _, t := range times {
		// First record can't be fixed, so skip processing.
		if prev != nil {
			if t.Time().Equal(prev.Time()) {
				prev.SetTime(prev.Time().Add(-time.Hour))
			}
		}
		prev = t
	}
}
