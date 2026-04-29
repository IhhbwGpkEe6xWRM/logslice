package parser

import (
	"errors"
	"strings"
	"time"
)

// LogLine represents a parsed log line with its timestamp and raw content.
type LogLine struct {
	Timestamp time.Time
	Raw       string
}

// ErrNoTimestamp is returned when a log line does not contain a parseable timestamp.
var ErrNoTimestamp = errors.New("no parseable timestamp found in line")

// commonPrefixFormats lists timestamp formats that typically appear at the
// start of a log line. They are tried in order, most specific first.
var commonPrefixFormats = []string{
	"2006-01-02T15:04:05.999999999Z07:00", // RFC3339Nano
	"2006-01-02T15:04:05Z07:00",           // RFC3339
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700", // Common Log Format
	"Jan  2 15:04:05",            // syslog
	"Jan 02 15:04:05",            // syslog variant
}

// ParseLine attempts to extract a timestamp from the beginning of the raw log
// line. It tries each known format and returns the first match.
func ParseLine(raw string) (LogLine, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return LogLine{}, ErrNoTimestamp
	}

	for _, format := range commonPrefixFormats {
		// Only look at the prefix of the line that is as long as the format.
		prefixLen := len(format)
		if prefixLen > len(trimmed) {
			prefixLen = len(trimmed)
		}

		for end := prefixLen; end >= 10; end-- {
			t, err := time.Parse(format, trimmed[:end])
			if err == nil {
				return LogLine{Timestamp: t, Raw: raw}, nil
			}
		}
	}

	return LogLine{}, ErrNoTimestamp
}

// InRange reports whether the log line's timestamp falls within [start, end]
// (both bounds inclusive). If start or end is zero they are treated as open
// bounds.
func (l LogLine) InRange(start, end time.Time) bool {
	if !start.IsZero() && l.Timestamp.Before(start) {
		return false
	}
	if !end.IsZero() && l.Timestamp.After(end) {
		return false
	}
	return true
}
