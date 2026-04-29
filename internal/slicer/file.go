package slicer

import (
	"fmt"
	"io"
	"os"
)

// SliceFile opens the file at srcPath, slices it according to opts, and writes
// results to dst. Returns the number of matching lines written.
func SliceFile(srcPath string, dst io.Writer, opts Options) (int, error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return 0, fmt.Errorf("slicer: open %q: %w", srcPath, err)
	}
	defer f.Close()

	s := New(opts)
	count, err := s.Slice(f, dst)
	if err != nil {
		return count, fmt.Errorf("slicer: slice %q: %w", srcPath, err)
	}
	return count, nil
}

// SliceFileToFile opens srcPath, slices it, and writes results to dstPath.
// The destination file is created or truncated.
func SliceFileToFile(srcPath, dstPath string, opts Options) (int, error) {
	out, err := os.Create(dstPath)
	if err != nil {
		return 0, fmt.Errorf("slicer: create %q: %w", dstPath, err)
	}
	defer out.Close()

	return SliceFile(srcPath, out, opts)
}
