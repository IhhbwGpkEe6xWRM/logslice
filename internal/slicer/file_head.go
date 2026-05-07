package slicer

import (
	"fmt"
	"io"
	"os"
	"time"
)

// HeadFile opens path and returns the first N lines/duration from [from, to].
func HeadFile(path string, from, to time.Time, opts HeadOptions) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	lines, err := SliceWithHead(f, from, to, opts)
	if err != nil {
		return nil, err
	}

	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = l.Raw
	}
	return out, nil
}

// HeadFileToWriter opens path, applies head options, and writes results to w.
func HeadFileToWriter(path string, from, to time.Time, opts HeadOptions, w io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	lines, err := SliceWithHead(f, from, to, opts)
	if err != nil {
		return err
	}

	for _, l := range lines {
		if _, err := fmt.Fprintln(w, l.Raw); err != nil {
			return err
		}
	}
	return nil
}
