package energiatili

import (
	"errors"
	"io"
	"strings"
)

// ErrorDateFormat is returned if there is no matching end parenthesis for
// timestamp
var ErrorDateFormat = errors.New("Error in date format")

const (
	startData = "var model = "
	endData   = ";"
	startDate = "new Date("
	endDate   = ")"
)

// dateConverter replaces "new Date(1234)" with "1234"
func dateConverter(input string, w io.Writer) (n int, err error) {
	r := strings.NewReplacer(startDate, "", endDate, "")
	n, err = r.WriteString(w, input)
	return n, err
}

// html2json finds an embedded javascript object in HTML and converts timestamps to integers
func html2json(s string, w io.Writer) (n int, err error) {
	// Find var model = ....
	start := strings.Index(s, startData)
	end := start + strings.Index(s[start:], endData)
	modelData := s[start+len(startData) : end]
	// Convert dates to JSON (ie. strip new `Date(` and `)`
	n, err = dateConverter(modelData, w)
	return n, err
}
