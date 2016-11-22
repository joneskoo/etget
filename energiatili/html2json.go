package energiatili

import (
	"errors"
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
func dateConverter(input string) (output string) {
	r := strings.NewReplacer(startDate, "", endDate, "")
	return r.Replace(input)
}

// html2json finds an embedded javascript object in HTML and converts timestamps to integers
func html2json(data []byte) (json string, err error) {
	s := string(data)
	// Find var model = ....
	start := strings.Index(s, startData)
	end := start + strings.Index(s[start:], endData)
	modelData := s[start+len(startData) : end]
	// Convert dates to JSON (ie. strip new `Date(` and `)`
	modelData = dateConverter(modelData)
	return modelData, nil
}
