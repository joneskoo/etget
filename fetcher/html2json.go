package fetcher

import (
	"bytes"
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
func dateConverter(input string) (output string, err error) {
	pos := 0
	var buffer bytes.Buffer
	for {
		// Look for new Date(
		start := strings.Index(input[pos:], startDate)
		if start == -1 {
			// No more matches, pass through the rest
			if _, err = buffer.WriteString(input[pos:]); err != nil {
				return "", err
			}
			break
		}
		start += pos
		// Look for closing )
		end := strings.Index(input[start:], endDate)
		if end == -1 {
			// Invalid format; got start marker but no end
			return "", ErrorDateFormat
		}
		end += start
		// Write everything except start and stop markers
		if _, err = buffer.WriteString(input[pos:start]); err != nil {
			return "", err
		}
		if _, err = buffer.WriteString(input[start+len(startDate) : end]); err != nil {
			return "", err
		}
		pos = end + 1
	}
	return buffer.String(), nil
}

// html2json finds an embedded javascript object in HTML and converts timestamps to integers
func html2json(data []byte) (json string, err error) {
	s := string(data)
	// Find var model = ....
	start := strings.Index(s, startData)
	end := start + strings.Index(s[start:], endData)
	modelData := s[start+len(startData) : end]
	// Convert dates to JSON
	modelData, err = dateConverter(modelData)
	return modelData, nil
}
