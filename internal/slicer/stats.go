package slicer

import (
	"fmt"
	"io"
	"time"
)

// Stats holds metrics collected during a slice operation.
type Stats struct {
	TotalLines   int
	MatchedLines int
	SkippedLines int
	UnparsedLines int
	Duration     time.Duration
	FirstMatch   *time.Time
	LastMatch    *time.Time
}

// Summary writes a human-readable summary of the stats to w.
func (s *Stats) Summary(w io.Writer) {
	fmt.Fprintf(w, "Total lines:    %d\n", s.TotalLines)
	fmt.Fprintf(w, "Matched lines:  %d\n", s.MatchedLines)
	fmt.Fprintf(w, "Skipped lines:  %d\n", s.SkippedLines)
	fmt.Fprintf(w, "Unparsed lines: %d\n", s.UnparsedLines)
	fmt.Fprintf(w, "Duration:       %s\n", s.Duration.Round(time.Millisecond))
	if s.FirstMatch != nil {
		fmt.Fprintf(w, "First match:    %s\n", s.FirstMatch.Format(time.RFC3339))
	}
	if s.LastMatch != nil {
		fmt.Fprintf(w, "Last match:     %s\n", s.LastMatch.Format(time.RFC3339))
	}
}

// MatchRate returns the fraction of total lines that matched, or 0 if no lines
// were processed.
func (s *Stats) MatchRate() float64 {
	if s.TotalLines == 0 {
		return 0
	}
	return float64(s.MatchedLines) / float64(s.TotalLines)
}
