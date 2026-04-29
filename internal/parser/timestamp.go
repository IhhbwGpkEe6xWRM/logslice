package parser

import (
	"fmt"
	"time"
)

// CommonFormats holds the timestamp formats logslice attempts to detect.
var CommonFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
}

// ParseTimestamp attempts to parse a raw string into a time.Time using
// the list of CommonFormats. Returns the parsed time and the matched
// format string, or an error if no format matched.
func ParseTimestamp(raw string) (time.Time, string, error) {
	for _, layout := range CommonFormats {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, layout, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("parser: unrecognized timestamp format: %q", raw)
}

// MustParseTimestamp is like ParseTimestamp but panics on error.
// Intended for use in tests and CLI flag validation.
func MustParseTimestamp(raw string) time.Time {
	t, _, err := ParseTimestamp(raw)
	if err != nil {
		panic(err)
	}
	return t
}
