package slicer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// RotateSliceOptions configures slicing across rotated log files.
type RotateSliceOptions struct {
	// Dir is the directory containing rotated log files.
	Dir string
	// Pattern is a glob pattern to match log files (e.g. "app.log*").
	Pattern string
	// From and To define the time range to extract.
	From time.Time
	To   time.Time
}

// RotateSlice reads all matching rotated log files in Dir, merges them
// chronologically, and writes lines within [From, To) to w.
func RotateSlice(opts RotateSliceOptions, w io.Writer) error {
	if opts.Dir == "" {
		return fmt.Errorf("rotate slice: dir must not be empty")
	}
	if opts.Pattern == "" {
		return fmt.Errorf("rotate slice: pattern must not be empty")
	}
	if !opts.To.IsZero() && opts.From.After(opts.To) {
		return fmt.Errorf("rotate slice: from %s is after to %s", opts.From, opts.To)
	}

	glob := filepath.Join(opts.Dir, opts.Pattern)
	matches, err := filepath.Glob(glob)
	if err != nil {
		return fmt.Errorf("rotate slice: glob %q: %w", glob, err)
	}
	if len(matches) == 0 {
		return fmt.Errorf("rotate slice: no files matched %q", glob)
	}

	sort.Strings(matches)

	var readers []io.Reader
	var closers []io.Closer
	defer func() {
		for _, c := range closers {
			_ = c.Close()
		}
	}()

	for _, path := range matches {
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("rotate slice: open %q: %w", path, err)
		}
		readers = append(readers, f)
		closers = append(closers, f)
	}

	return MergeSlice(readers, opts.From, opts.To, w)
}

// ListRotatedFiles returns the sorted list of file paths that match the
// pattern inside dir, along with the detected format of the first file.
func ListRotatedFiles(dir, pattern string) ([]string, string, error) {
	glob := filepath.Join(dir, pattern)
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, "", fmt.Errorf("list rotated files: %w", err)
	}
	sort.Strings(matches)

	format := "unknown"
	if len(matches) > 0 {
		f, err := os.Open(matches[0])
		if err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			var lines []string
			for scanner.Scan() && len(lines) < 5 {
				lines = append(lines, scanner.Text())
			}
			if detected := parser.DetectFormat(strings.NewReader(strings.Join(lines, "\n"))); detected != "" {
				format = detected
			}
		}
	}
	return matches, format, nil
}
