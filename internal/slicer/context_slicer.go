package slicer

import (
	"bufio"
	"io"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// ContextOptions controls how many lines before/after a matching line are included.
type ContextOptions struct {
	Before int
	After  int
}

// SliceWithContext extracts lines in [from, to] and includes up to Before lines
// preceding each match and After lines following each match.
func SliceWithContext(r io.Reader, from, to time.Time, opts ContextOptions) ([]parser.LogLine, error) {
	scanner := bufio.NewScanner(r)

	var all []parser.LogLine
	for scanner.Scan() {
		line := parser.ParseLine(scanner.Text())
		all = append(all, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(all) == 0 {
		return nil, nil
	}

	// Mark which indices are primary matches.
	matched := make([]bool, len(all))
	for i, l := range all {
		if l.InRange(from, to) {
			matched[i] = true
		}
	}

	// Expand with context window.
	included := make([]bool, len(all))
	for i, m := range matched {
		if !m {
			continue
		}
		start := i - opts.Before
		if start < 0 {
			start = 0
		}
		end := i + opts.After
		if end >= len(all) {
			end = len(all) - 1
		}
		for j := start; j <= end; j++ {
			included[j] = true
		}
	}

	var result []parser.LogLine
	for i, l := range all {
		if included[i] {
			result = append(result, l)
		}
	}
	return result, nil
}
