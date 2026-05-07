package slicer

import (
	"io"
	"os"
	"time"
)

// GrepFile opens path and writes matching lines in [from, to] to w.
func GrepFile(path string, from, to time.Time, opts GrepOptions, w io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return GrepSlice(f, from, to, opts, w)
}

// GrepFileToFile opens src, greps it, and writes results to dst (created/truncated).
func GrepFileToFile(src, dst string, from, to time.Time, opts GrepOptions) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	return GrepFile(src, from, to, opts, out)
}
