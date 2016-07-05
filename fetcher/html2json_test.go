package fetcher

import (
	"testing"
)

func TestDateConverterPositive(t *testing.T) {
	// Positive tests
	cases := []struct {
		in,
		want string
		error bool
	}{
		{in: "no date in string", want: "no date in string"},
		{in: "no (date in string", want: "no (date in string"},
		{in: "no )date in string", want: "no )date in string"},
		{in: "new Date(123)", want: "\"123\""},
		{in: "asdf, new Date(123), new Date(124), foo", want: "asdf, \"123\", \"124\", foo"},
		{in: "a adf new Date(123) foo", want: "a adf \"123\" foo"},
		{in: "\"StartValue\":new Date(1325368800000),\"", want: "\"StartValue\":\"1325368800000\",\""},
		{in: "\"StartValue\":new Date(1325368800000,\"", error: true},
	}

	for _, c := range cases {
		out, err := dateConverter(c.in)
		if c.error && err == nil {
			t.Errorf("dateConverter(%q); expected error, got err=%q", c.in, err)
		}
		if !c.error && err != nil {
			t.Errorf("dateConverter(%q); expected err=nil, got err=%q", c.in, err)
		}
		if out != c.want {
			t.Errorf("dateConverter(%q); expected %q, got %q", c.in, c.want, out)
		}
	}
}
