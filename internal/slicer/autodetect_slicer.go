package slicer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"time"

	"logslice/internal/parser"
)

// AutoSliceOptions configures the auto-detecting slicer.
type AutoSliceOptions struct {
	From       time.Time
	To         time.Time
	SampleSize int // number of lines to probe for format detection
}

// AutoSlice reads from r, auto-detects the timestamp format from the first
// SampleSize lines, then extracts log lines whose timestamps fall within
// [From, To] and writes them to w.
func AutoSlice(r io.Reader, w io.Writer, opts AutoSliceOptions) error {
	if opts.SampleSize <= 0 {
		opts.SampleSize = 20
	}

	// Buffer the entire input so we can probe then re-read.
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("autodetect slicer: read: %w", err)
	}

	sample := parser.SampleReader(bytes.NewReader(data), opts.SampleSize)
	if !sample.Detected {
		return fmt.Errorf("autodetect slicer: could not detect timestamp format")
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		ll, ok := parser.ParseLine(line)
		if !ok {
			continue
		}
		if ll.InRange(opts.From, opts.To) {
			if _, err := fmt.Fprintln(w, line); err != nil {
				return fmt.Errorf("autodetect slicer: write: %w", err)
			}
		}
	}
	return scanner.Err()
}

// DetectedFormat returns the format string detected from the first lines of r
// without consuming the reader for further use.
func DetectedFormat(r io.Reader, sampleSize int) (string, error) {
	if sampleSize <= 0 {
		sampleSize = 20
	}
	result := parser.SampleReader(r, sampleSize)
	if !result.Detected {
		return "", fmt.Errorf("could not detect timestamp format")
	}
	return result.Format, nil
}
