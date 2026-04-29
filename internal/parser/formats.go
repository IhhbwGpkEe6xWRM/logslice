// Package parser provides utilities for parsing structured log lines,
// including timestamp extraction from a variety of common log formats.
package parser

import "time"

// TimestampFormat describes a named timestamp layout for use with time.Parse.
type TimestampFormat struct {
	// Name is a human-readable identifier for the format.
	Name string
	// Layout is the Go time.Parse reference layout string.
	Layout string
	// PrefixLen is the number of characters to read from the start of a log
	// line when attempting to parse the timestamp. Zero means use the full
	// layout length heuristic.
	PrefixLen int
}

// KnownFormats is the ordered list of timestamp formats that logslice
// recognises. Formats are tried in order; the first successful parse wins.
//
// When adding a new format, place more-specific (longer) layouts before
// shorter or ambiguous ones to avoid false positives.
var KnownFormats = []TimestampFormat{
	{
		Name:      "RFC3339Nano",
		Layout:    time.RFC3339Nano,
		PrefixLen: 35, // e.g. "2006-01-02T15:04:05.999999999Z07:00"
	},
	{
		Name:      "RFC3339",
		Layout:    time.RFC3339,
		PrefixLen: 25, // e.g. "2006-01-02T15:04:05Z07:00"
	},
	{
		Name:      "DateTime",
		Layout:    "2006-01-02 15:04:05",
		PrefixLen: 19,
	},
	{
		Name:      "DateTimeMilli",
		Layout:    "2006-01-02 15:04:05.000",
		PrefixLen: 23,
	},
	{
		Name:      "DateTimeMicro",
		Layout:    "2006-01-02 15:04:05.000000",
		PrefixLen: 26,
	},
	{
		Name:      "ApacheCLF",
		Layout:    "02/Jan/2006:15:04:05 -0700",
		PrefixLen: 26,
	},
	{
		Name:      "Syslog",
		Layout:    "Jan _2 15:04:05",
		PrefixLen: 15,
	},
	{
		Name:      "UnixDate",
		Layout:    time.UnixDate,
		PrefixLen: 28,
	},
}

// TryParseTimestampPrefix attempts to parse a timestamp from the leading
// characters of s using each entry in KnownFormats. It returns the parsed
// time and the remaining suffix of s after the timestamp, or a zero time and
// the original string if no format matched.
func TryParseTimestampPrefix(s string) (t time.Time, rest string, ok bool) {
	for _, f := range KnownFormats {
		prefixLen := f.PrefixLen
		if prefixLen <= 0 || prefixLen > len(s) {
			prefixLen = len(s)
		}
		candidate := s[:prefixLen]
		parsed, err := time.Parse(f.Layout, candidate)
		if err == nil {
			return parsed, s[prefixLen:], true
		}
		// Also try the full remaining string for formats where the prefix
		// length is just a hint and the actual token may be shorter.
		if prefixLen < len(s) {
			parsed, err = time.Parse(f.Layout, s)
			if err == nil {
				return parsed, "", true
			}
		}
	}
	return time.Time{}, s, false
}
