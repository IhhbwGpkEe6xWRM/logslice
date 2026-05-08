package slicer

import (
	"io"
	"os"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// ContextFile opens the named log file and extracts lines in [from, to] with
// surrounding context lines. Results are written to w.
func ContextFile(path string, from, to time.Time, opts ContextOptions, w io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	lines, err := SliceWithContext(f, from, to, opts)
	if err != nil {
		return err
	}
	return writeLines(w, lines)
}

// ContextFileToFile extracts context lines from src and writes them to dst.
func ContextFileToFile(src, dst string, from, to time.Time, opts ContextOptions) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	return ContextFile(src, from, to, opts, out)
}

// writeLines writes each LogLine's raw text followed by a newline to w.
func writeLines(w io.Writer, lines []parser.LogLine) error {
	for _, l := range lines {
		if _, err := io.WriteString(w, l.Raw+"\n"); err != nil {
			return err
		}
	}
	return nil
}
