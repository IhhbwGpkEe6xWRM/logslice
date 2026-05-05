package slicer

import (
	"fmt"
	"io"
	"os"
)

// TailFile opens the named log file and returns a TailResult according to
// opts. The file is closed before returning.
func TailFile(path string, opts TailOptions) (TailResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return TailResult{}, fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()
	return SliceWithTail(f, opts)
}

// TailFileToWriter opens the named log file, extracts the tail window, and
// writes each matched line followed by a newline to w.
func TailFileToWriter(path string, opts TailOptions, w io.Writer) (Stats, error) {
	res, err := TailFile(path, opts)
	if err != nil {
		return Stats{}, err
	}
	for _, line := range res.Lines {
		if _, werr := fmt.Fprintln(w, line); werr != nil {
			return res.Stats, fmt.Errorf("write: %w", werr)
		}
	}
	return res.Stats, nil
}
