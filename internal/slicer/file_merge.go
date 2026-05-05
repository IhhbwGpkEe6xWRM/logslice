package slicer

import (
	"fmt"
	"io"
	"os"
	"time"
)

// MergeFiles opens each path in paths, merges their log lines in chronological
// order within [from, to], and writes results to w.
func MergeFiles(paths []string, from, to time.Time, w io.Writer) (int, error) {
	if len(paths) == 0 {
		return 0, fmt.Errorf("merge: no input files provided")
	}

	files := make([]*os.File, 0, len(paths))
	readers := make([]interface{ Read([]byte) (int, error) }, 0, len(paths))

	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			for _, open := range files {
				open.Close()
			}
			return 0, fmt.Errorf("merge: open %q: %w", p, err)
		}
		files = append(files, f)
		readers = append(readers, f)
	}
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()

	return MergeSlice(readers, from, to, w)
}

// MergeFilesToFile merges log lines from paths into the file at outPath.
func MergeFilesToFile(paths []string, from, to time.Time, outPath string) (int, error) {
	out, err := os.Create(outPath)
	if err != nil {
		return 0, fmt.Errorf("merge: create output %q: %w", outPath, err)
	}
	defer out.Close()
	return MergeFiles(paths, from, to, out)
}
