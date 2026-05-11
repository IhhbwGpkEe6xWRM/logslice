package slicer

import (
	"io"
	"os"
	"time"
)

// AnnotateFile opens path and writes annotated lines within [from, to] to w.
func AnnotateFile(path string, from, to time.Time, opts AnnotateOptions, w io.Writer) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return SliceAnnotated(f, from, to, opts, w)
}

// AnnotateFileToFile opens src, annotates lines within [from, to], and writes
// results to dst, creating or truncating it.
func AnnotateFileToFile(src, dst string, from, to time.Time, opts AnnotateOptions) (int, error) {
	in, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	return SliceAnnotated(in, from, to, opts, out)
}
