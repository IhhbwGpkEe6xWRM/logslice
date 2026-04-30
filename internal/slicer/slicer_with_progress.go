package slicer

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/user/logslice/internal/parser"
)

// SliceWithProgress slices log lines within [from, to] and writes matching
// lines to dst, reporting progress via reporter.
// The reporter may be nil, in which case no progress is reported.
func SliceWithProgress(
	src io.Reader,
	dst io.Writer,
	from, to time.Time,
	reporter *ProgressReporter,
) (Stats, error) {
	scanner := bufio.NewScanner(src)
	collector := newStatsCollector()

	for scanner.Scan() {
		raw := scanner.Text()
		line := parser.ParseLine(raw)
		byteLen := len(raw) + 1 // +1 for newline

		matched := false
		if line.InRange(from, to) {
			if _, err := fmt.Fprintln(dst, raw); err != nil {
				return collector.Stats(), fmt.Errorf("writing output: %w", err)
			}
			matched = true
		}

		collector.record(line, from, to)
		if reporter != nil {
			reporter.RecordLine(byteLen, matched)
		}
	}

	if err := scanner.Err(); err != nil {
		return collector.Stats(), fmt.Errorf("reading input: %w", err)
	}

	if reporter != nil {
		reporter.Finish()
	}
	return collector.Stats(), nil
}
