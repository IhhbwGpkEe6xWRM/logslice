package slicer

import (
	"bufio"
	"io"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// SliceWithStats performs the same operation as the Slicer but also collects
// and returns Stats about the operation.
func (s *Slicer) SliceWithStats(r io.Reader, w io.Writer) (*Stats, error) {
	start := time.Now()
	stats := &Stats{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		stats.TotalLines++

		ll, err := parser.ParseLine(line)
		if err != nil || ll.Timestamp.IsZero() {
			stats.UnparsedLines++
			continue
		}

		if ll.Timestamp.Before(s.from) {
			stats.SkippedLines++
			continue
		}
		if ll.Timestamp.After(s.to) {
			stats.SkippedLines++
			continue
		}

		if _, err := fmt.Fprintln(w, line); err != nil {
			return stats, err
		}

		stats.MatchedLines++
		t := ll.Timestamp
		if stats.FirstMatch == nil {
			stats.FirstMatch = &t
		}
		stats.LastMatch = &t
	}

	if err := scanner.Err(); err != nil {
		return stats, err
	}

	stats.Duration = time.Since(start)
	return stats, nil
}
