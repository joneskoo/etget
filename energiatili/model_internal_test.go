package energiatili

import "testing"
import "reflect"
import "time"

func TestTrimTrailingZeros(t *testing.T) {
	date := time.Date(2016, 12, 11, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		in   []float64
		want []float64
	}{
		{
			in:   []float64{},
			want: []float64{},
		},
		{
			in:   []float64{0},
			want: []float64{},
		},
		{
			in:   []float64{1},
			want: []float64{1},
		},
		{
			in:   []float64{1, 0, 0},
			want: []float64{1},
		},
		{
			in:   []float64{1, 0, 0, 1},
			want: []float64{1, 0, 0, 1},
		},
	}
	for _, c := range cases {
		records := []Record{}
		for _, v := range c.in {
			records = append(records, Record{Timestamp: date, Value: v})
		}

		got := []float64{}
		for _, v := range trimTrailingZeros(records) {
			got = append(got, v.Value)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("want trimTrailingZeros(%#v) = %#v, got %#v", c.in, c.want, got)
		}
	}
}
