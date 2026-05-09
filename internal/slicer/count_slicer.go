package slicer

import (
	"bufio"
	"io"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// CountResult holds the result of a count operation over a log range.
type CountResult struct {
	Total    int
	Matched  int
	ByLevel  map[string]int
	ByMinute map[string]int
}

// SliceCount scans lines from r within [from, to] and returns aggregated counts.
// It does not write output — it only counts matching log lines.
func SliceCount(r io.Reader, from, to time.Time) (*CountResult, error) {
	result := &CountResult{
		ByLevel:  make(map[string]int),
		ByMinute: make(map[string]int),
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		raw := scanner.Text()
		if raw == "" {
			continue
		}
		result.Total++

		line, err := parser.ParseLine(raw)
		if err != nil || line.Timestamp.IsZero() {
			continue
		}

		if !line.InRange(from, to) {
			continue
		}

		result.Matched++

		if line.Level != "" {
			result.ByLevel[line.Level]++
		}

		minKey := line.Timestamp.UTC().Format("2006-01-02T15:04")
		result.ByMinute[minKey]++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
