package energiatili_test

import (
	"testing"
	"time"

	"github.com/joneskoo/etget/energiatili"
)

func must(t time.Time, err error) time.Time {
	if err != nil {
		panic(err)
	}
	return t
}

type step struct {
	in   float64
	want time.Time
}

type testcase []step

func TestTimeFixer(t *testing.T) {
	helsinki, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		panic(err)
	}
	cases := []testcase{
		// Winter time
		{
			step{1325379600000.0, must(time.ParseInLocation(time.UnixDate, "Sun Jan  1 01:00:00 EET 2012", helsinki))},
			step{1325383200000.0, must(time.ParseInLocation(time.UnixDate, "Sun Jan  1 02:00:00 EET 2012", helsinki))},
			step{1325386800000.0, must(time.ParseInLocation(time.UnixDate, "Sun Jan  1 03:00:00 EET 2012", helsinki))},
			step{1325390400000.0, must(time.ParseInLocation(time.UnixDate, "Sun Jan  1 04:00:00 EET 2012", helsinki))},
		},
		// Summer time
		{
			step{1467342000000.0, must(time.ParseInLocation(time.UnixDate, "Fri Jul  1 00:00:00 UTC 2016", helsinki))},
			step{1467345600000.0, must(time.ParseInLocation(time.UnixDate, "Fri Jul  1 01:00:00 UTC 2016", helsinki))},
			step{1467349200000.0, must(time.ParseInLocation(time.UnixDate, "Fri Jul  1 02:00:00 UTC 2016", helsinki))},
			step{1467352800000.0, must(time.ParseInLocation(time.UnixDate, "Fri Jul  1 03:00:00 UTC 2016", helsinki))},
		},
		// Winter -> Summer
		{
			step{1459040400000.0, must(time.ParseInLocation(time.UnixDate, "Sun Mar 26 23:00:00 UTC 2016", helsinki))},
			step{1459044000000.0, must(time.ParseInLocation(time.UnixDate, "Sun Mar 27 00:00:00 UTC 2016", helsinki))},
			step{1459051200000.0, must(time.ParseInLocation(time.UnixDate, "Sun Mar 27 01:00:00 UTC 2016", helsinki))},
			step{1459054800000.0, must(time.ParseInLocation(time.UnixDate, "Sun Mar 27 02:00:00 UTC 2016", helsinki))},
			step{1459058400000.0, must(time.ParseInLocation(time.UnixDate, "Sun Mar 27 03:00:00 UTC 2016", helsinki))},
		},
		// Summer -> Winter
		{
			step{1445734800000.0, must(time.ParseInLocation(time.UnixDate, "Sun Oct 24 22:00:00 UTC 2015", helsinki))},
			step{1445738400000.0, must(time.ParseInLocation(time.UnixDate, "Sun Oct 24 23:00:00 UTC 2015", helsinki))},
			step{1445742000000.0, must(time.ParseInLocation(time.UnixDate, "Sun Oct 25 00:00:00 UTC 2015", helsinki))},
			step{1445742000000.0, must(time.ParseInLocation(time.UnixDate, "Sun Oct 25 01:00:00 UTC 2015", helsinki))},
			step{1445745600000.0, must(time.ParseInLocation(time.UnixDate, "Sun Oct 25 02:00:00 UTC 2015", helsinki))},
			step{1445749200000.0, must(time.ParseInLocation(time.UnixDate, "Sun Oct 25 03:00:00 UTC 2015", helsinki))},
		},
	}

	for _, tc := range cases {
		var fixer energiatili.TimeFixer
		for _, step := range tc {
			ts, err := fixer.ParseBrokenTime(step.in)
			if err != nil {
				t.Error(err)
			}
			want := step.want
			if !ts.Equal(want) {
				t.Errorf("Fail, Got %v != Want %v, diff=%v", ts.Local(), want.Local(), ts.Sub(want))
			}
		}

	}
}
