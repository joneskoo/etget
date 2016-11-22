package energiatili

import (
	"bytes"
	"testing"
)

func TestDateConverterPositive(t *testing.T) {
	// Positive tests
	cases := []struct {
		in,
		want string
	}{
		{in: "no date in string", want: "no date in string"},
		{in: "no (date in string", want: "no (date in string"},
		{in: "no )date in string", want: "no date in string"},
		{in: "new Date(123)", want: "123"},
		{in: "asdf, new Date(123), new Date(124), foo", want: "asdf, 123, 124, foo"},
		{in: "a adf new Date(123) foo", want: "a adf 123 foo"},
		{in: "\"StartValue\":new Date(1325368800000),\"", want: "\"StartValue\":1325368800000,\""},
	}

	for _, c := range cases {
		buf := new(bytes.Buffer)
		n, err := dateConverter(c.in, buf)
		if n != buf.Len() {
			t.Errorf("n, err := dateConverter(%q, w); expected n (%d) to equal bytes written (%d)", c.in, n, buf.Len())
		}
		if err != nil {
			t.Errorf("n, err := dateConverter(%q, w); got error: %s", c.in, err)
		}
		if buf.String() != c.want {
			t.Errorf("dateConverter(%q, w); expected %q, got %q", c.in, c.want, buf.String())
		}
	}
}
