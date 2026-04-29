package slicer

import (
	"bufio"
	"io"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// Options configures the slicing behavior.
type Options struct {
	From      time.Time
	To        time.Time
	Inclusive bool
}

// Slicer extracts log lines within a time range from a reader.
type Slicer struct {
	opts Options
}

// New creates a new Slicer with the given options.
func New(opts Options) *Slicer {
	return &Slicer{opts: opts}
}

// Slice reads from src, writes matching lines to dst, and returns the count
// of lines written and any error encountered.
func (s *Slicer) Slice(src io.Reader, dst io.Writer) (int, error) {
	scanner := bufio.NewScanner(src)
	writer := bufio.NewWriter(dst)
	defer writer.Flush()

	count := 0
	for scanner.Scan() {
		raw := scanner.Text()
		line := parser.ParseLine(raw)

		if line.InRange(s.opts.From, s.opts.To) {
			if _, err := writer.WriteString(raw + "\n"); err != nil {
				return count, err
			}
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		return count, err
	}

	return count, nil
}
