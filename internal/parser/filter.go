package parser

import (
	"regexp"
	"strings"
	"time"
)

// LineFilter defines optional criteria for filtering log lines beyond time range.
type LineFilter struct {
	// Level filters lines by log level (e.g. "ERROR", "WARN"). Empty means no filter.
	Level string
	// Pattern filters lines matching a regular expression. Nil means no filter.
	Pattern *regexp.Regexp
}

// NewLineFilter constructs a LineFilter from raw string inputs.
// level is matched case-insensitively. pattern is compiled as a regexp.
// Returns an error if pattern is not a valid regexp.
func NewLineFilter(level, pattern string) (*LineFilter, error) {
	f := &LineFilter{
		Level: strings.ToUpper(strings.TrimSpace(level)),
	}
	if pattern != "" {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		f.Pattern = re
	}
	return f, nil
}

// Match reports whether the given raw log line satisfies the filter.
// ts is the parsed timestamp of the line (may be zero if unparsed).
func (f *LineFilter) Match(line string, _ time.Time) bool {
	if f == nil {
		return true
	}
	if f.Level != "" {
		if !strings.Contains(strings.ToUpper(line), f.Level) {
			return false
		}
	}
	if f.Pattern != nil {
		if !f.Pattern.MatchString(line) {
			return false
		}
	}
	return true
}
