package slicer

import (
	"bufio"
	"io"
	"regexp"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// GrepOptions controls how GrepSlice matches and emits lines.
type GrepOptions struct {
	Pattern     *regexp.Regexp
	Invert      bool // emit lines that do NOT match the pattern
	Context     int  // number of surrounding lines to include
}

// GrepSlice extracts lines within [from, to] whose body matches opts.Pattern,
// optionally including Context lines before and after each match.
func GrepSlice(r io.Reader, from, to time.Time, opts GrepOptions, w io.Writer) error {
	if opts.Pattern == nil {
		return SliceWithProgress(r, from, to, nil, w)
	}

	scanner := bufio.NewScanner(r)
	var lines []parser.LogLine
	for scanner.Scan() {
		ll := parser.ParseLine(scanner.Text())
		lines = append(lines, ll)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Collect indices of lines that are in range and match the grep condition.
	var matchIdx []int
	for i, ll := range lines {
		if !ll.InRange(from, to) {
			continue
		}
		matched := opts.Pattern.MatchString(ll.Raw)
		if matched != opts.Invert {
			matchIdx = append(matchIdx, i)
		}
	}

	// Expand each match with context lines, deduplicating via a set.
	emit := make(map[int]struct{})
	for _, idx := range matchIdx {
		start := idx - opts.Context
		if start < 0 {
			start = 0
		}
		end := idx + opts.Context
		if end >= len(lines) {
			end = len(lines) - 1
		}
		for k := start; k <= end; k++ {
			emit[k] = struct{}{}
		}
	}

	bw := bufio.NewWriter(w)
	for i := 0; i < len(lines); i++ {
		if _, ok := emit[i]; !ok {
			continue
		}
		if _, err := bw.WriteString(lines[i].Raw + "\n"); err != nil {
			return err
		}
	}
	return bw.Flush()
}
