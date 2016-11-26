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

// FixDST fixes DST ambiguoity in a slice of values.
// All values must be hours with no minutes, seconds or nanoseconds part
// in sequential order.
func FixDST(data Interface) {
	for i := 1; i < data.Len(); i++ {
		prev := data.Time(i - 1)
		if data.Time(i).Equal(prev) {
			data.SetTime(i-1, prev.Add(-time.Hour))
		}
	}
}

// Interface must be implemented by values used with FixDST.
type Interface interface {
	// Len is the number of elements in the collection.
	Len() int

	// Time retrieves the timestamp value to be fixed.
	Time(i int) time.Time

	// SetTime sets the timestamp value to be fixed.
	SetTime(i int, t time.Time)
}

// Times is a slice of time.Time values and implements notz.Interface.
type Times []time.Time

func (t Times) Len() int                     { return len(t) }
func (t Times) Time(i int) time.Time         { return t[i] }
func (t Times) SetTime(i int, new time.Time) { t[i] = new }
