package parser

import (
	"strings"
	"time"
)

// knownFormats lists timestamp formats tried in order when parsing log prefixes.
// Each entry pairs a layout string with the expected character length of the
// timestamp portion so we can slice the line efficiently.
var knownFormats = []struct {
	layout string
	minLen int
}{
	// RFC3339 with nanoseconds
	{time.RFC3339Nano, 20},
	// RFC3339 basic (e.g. 2006-01-02T15:04:05Z)
	{time.RFC3339, 20},
	// Space-separated datetime with timezone offset
	{"2006-01-02 15:04:05 -0700", 25},
	// Space-separated datetime (assumed UTC)
	{"2006-01-02 15:04:05", 19},
	// Date only
	{"2006-01-02", 10},
}

// TryParseTimestampPrefix attempts to parse a timestamp from the beginning of
// line. It returns the parsed time, the remainder of the line after the
// timestamp (and any separating whitespace), and whether parsing succeeded.
func TryParseTimestampPrefix(line string) (time.Time, string, bool) {
	if len(line) == 0 {
		return time.Time{}, "", false
	}

	for _, f := range knownFormats {
		if len(line) < f.minLen {
			continue
		}

		// Try progressively longer prefixes starting at minLen to handle
		// variable-length suffixes like fractional seconds or timezone offsets.
		maxLen := f.minLen + 15
		if maxLen > len(line) {
			maxLen = len(line)
		}

		for end := maxLen; end >= f.minLen; end-- {
			candidate := line[:end]
			ts, err := time.Parse(f.layout, candidate)
			if err != nil {
				continue
			}
			rest := line[end:]
			rest = strings.TrimLeft(rest, " \t")
			return ts, rest, true
		}
	}

	return time.Time{}, "", false
}
