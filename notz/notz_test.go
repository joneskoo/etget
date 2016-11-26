package notz_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/joneskoo/etget/notz"
)

func TestTimeFixer(t *testing.T) {
	helsinki, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		panic(err)
	}
	cases := [][]time.Time{
		// Winter time
		{
			must(time.ParseInLocation(time.UnixDate, "Sun Jan  1 01:00:00 EET 2012", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Sun Jan  1 02:00:00 EET 2012", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Sun Jan  1 03:00:00 EET 2012", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Sun Jan  1 04:00:00 EET 2012", helsinki)),
		},
		// Summer time
		{
			must(time.ParseInLocation(time.UnixDate, "Fri Jul  1 00:00:00 UTC 2016", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Fri Jul  1 01:00:00 UTC 2016", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Fri Jul  1 02:00:00 UTC 2016", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Fri Jul  1 03:00:00 UTC 2016", helsinki)),
		},
		// Winter -> Summer
		{
			must(time.ParseInLocation(time.UnixDate, "Sun Mar 27 01:00:00 EET 2016", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Sun Mar 27 02:00:00 EET 2016", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Sun Mar 27 04:00:00 EEST 2016", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Sun Mar 27 05:00:00 EEST 2016", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Sun Mar 27 06:00:00 EEST 2016", helsinki)),
		},
		// Summer -> Winter
		{
			must(time.ParseInLocation(time.UnixDate, "Sun Oct 25 01:00:00 EEST 2015", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Sun Oct 25 02:00:00 EEST 2015", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Sun Oct 25 03:00:00 EEST 2015", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Sun Oct 25 03:00:00 EET 2015", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Sun Oct 25 04:00:00 EET 2015", helsinki)),
			must(time.ParseInLocation(time.UnixDate, "Sun Oct 25 05:00:00 EET 2015", helsinki)),
		},
	}

	for tc, testcase := range cases {
		// Create a time-zone unaware version (DST confusion)
		var times []notz.TimeSetter
		for _, tt := range testcase {
			tHelsinki := tt.In(helsinki)
			year, _, day := tHelsinki.Date()
			month := tHelsinki.Month()
			hour, min, sec := tHelsinki.Clock()
			nano := tHelsinki.Nanosecond()
			times = append(times, &mockTimeSetter{time.Date(year, month, day, hour, min, sec, nano, helsinki)})
		}
		notz.FixDST(times)

		for i, tt := range testcase {
			if !tt.Equal(times[i].Time()) {
				t.Errorf("testcase #%d[%d]: want %s, got %s", tc, i, tt.UTC(), times[i].Time().UTC())
			}
		}

	}
}

func ExampleFixDST() {
	helsinki, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		panic(err)
	}

	times := []notz.TimeSetter{
		&mockTimeSetter{time.Date(2015, 10, 25, 2, 0, 0, 0, helsinki)},
		&mockTimeSetter{time.Date(2015, 10, 25, 3, 0, 0, 0, helsinki)},
		&mockTimeSetter{time.Date(2015, 10, 25, 3, 0, 0, 0, helsinki)},
		&mockTimeSetter{time.Date(2015, 10, 25, 4, 0, 0, 0, helsinki)},
	}
	notz.FixDST(times)
	for _, t := range times {
		fmt.Printf("%s\n", t.Time())
	}
	// Output:
	// 2015-10-25 02:00:00 +0300 EEST
	// 2015-10-25 03:00:00 +0300 EEST
	// 2015-10-25 03:00:00 +0200 EET
	// 2015-10-25 04:00:00 +0200 EET
}

type mockTimeSetter struct {
	t time.Time
}

func (m mockTimeSetter) Time() time.Time      { return m.t }
func (m *mockTimeSetter) SetTime(t time.Time) { m.t = t }

func must(t time.Time, err error) time.Time {
	if err != nil {
		panic(err)
	}
	return t
}
