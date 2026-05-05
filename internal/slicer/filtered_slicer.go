package slicer

import (
	"bufio"
	"io"
	"time"

	"github.com/user/logslice/internal/parser"
)

// FilteredSliceOptions holds configuration for a filtered slice operation.
type FilteredSliceOptions struct {
	From   time.Time
	To     time.Time
	Filter *parser.LineFilter
}

// FilteredSliceResult holds the output of a filtered slice operation.
type FilteredSliceResult struct {
	LinesScanned  int
	LinesMatched  int
	LinesFiltered int
}

// SliceWithFilter reads lines from r, writes those in [from, to] that also
// satisfy filter to w. Returns a result summary.
func SliceWithFilter(r io.Reader, w io.Writer, opts FilteredSliceOptions) (FilteredSliceResult, error) {
	var res FilteredSliceResult
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		res.LinesScanned++

		ll, ok := parser.ParseLine(line)
		if !ok {
			continue
		}
		if !ll.InRange(opts.From, opts.To) {
			continue
		}
		res.LinesMatched++

		if !opts.Filter.Match(line, ll.Timestamp) {
			res.LinesFiltered++
			continue
		}

		if _, err := io.WriteString(w, line+"\n"); err != nil {
			return res, err
		}
	}
	return res, scanner.Err()
}
