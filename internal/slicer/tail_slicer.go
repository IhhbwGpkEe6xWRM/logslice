package slicer

import (
	"bufio"
	"io"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// TailOptions configures the tail slicing behaviour.
type TailOptions struct {
	Lines    int           // maximum number of lines to return (0 = unlimited)
	Duration time.Duration // return lines within the last Duration (0 = disabled)
}

// TailResult holds the output of a tail slice operation.
type TailResult struct {
	Lines []string
	Stats Stats
}

// SliceWithTail reads all lines from r that fall within the tail window
// defined by opts. The window is anchored to the timestamp of the last
// matched line, so the result is stable even when Duration and Lines are
// both set.
func SliceWithTail(r io.Reader, opts TailOptions) (TailResult, error) {
	scanner := bufio.NewScanner(r)

	var all []parser.LogLine
	for scanner.Scan() {
		text := scanner.Text()
		ll := parser.ParseLine(text)
		all = append(all, ll)
	}
	if err := scanner.Err(); err != nil {
		return TailResult{}, err
	}

	matched := filterTail(all, opts)

	var out []string
	for _, l := range matched {
		out = append(out, l.Raw)
	}

	s := Stats{
		TotalLines:   len(all),
		MatchedLines: len(matched),
	}
	return TailResult{Lines: out, Stats: s}, nil
}

func filterTail(lines []parser.LogLine, opts TailOptions) []parser.LogLine {
	if len(lines) == 0 {
		return nil
	}

	// Find the latest timestamp in the set.
	var latest time.Time
	for _, l := range lines {
		if l.Timestamp != nil && l.Timestamp.After(latest) {
			latest = *l.Timestamp
		}
	}

	var candidates []parser.LogLine
	for _, l := range lines {
		if opts.Duration > 0 && l.Timestamp != nil {
			cutoff := latest.Add(-opts.Duration)
			if l.Timestamp.Before(cutoff) {
				continue
			}
		}
		candidates = append(candidates, l)
	}

	if opts.Lines > 0 && len(candidates) > opts.Lines {
		candidates = candidates[len(candidates)-opts.Lines:]
	}
	return candidates
}
