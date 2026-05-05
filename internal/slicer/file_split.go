package slicer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SplitOptions controls how a log file is split into time-based chunks.
type SplitOptions struct {
	// BucketDuration is the width of each output chunk (e.g. time.Hour for hourly files).
	BucketDuration time.Duration

	// OutputDir is the directory where split files will be written.
	// Defaults to the directory of the input file when empty.
	OutputDir string

	// Prefix is prepended to each output filename.
	// When empty the base name of the input file (without extension) is used.
	Prefix string

	// Extension is appended to each output filename, e.g. ".log".
	// When empty the extension of the input file is preserved.
	Extension string
}

// SplitResult describes a single chunk produced by SplitFile.
type SplitResult struct {
	Path      string
	From      time.Time
	To        time.Time
	LineCount int
}

// SplitFile reads inputPath, groups log lines into fixed-duration buckets and
// writes each bucket to a separate file inside opts.OutputDir.
// It returns one SplitResult per non-empty bucket.
func SplitFile(inputPath string, from, to time.Time, opts SplitOptions) ([]SplitResult, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", inputPath, err)
	}
	defer f.Close()

	if opts.BucketDuration <= 0 {
		return nil, fmt.Errorf("BucketDuration must be positive")
	}

	// Resolve output directory.
	outDir := opts.OutputDir
	if outDir == "" {
		outDir = filepath.Dir(inputPath)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, fmt.Errorf("create output dir %s: %w", outDir, err)
	}

	// Resolve filename prefix and extension.
	base := filepath.Base(inputPath)
	prefix := opts.Prefix
	if prefix == "" {
		prefix = strings.TrimSuffix(base, filepath.Ext(base))
	}
	ext := opts.Extension
	if ext == "" {
		ext = filepath.Ext(base)
		if ext == "" {
			ext = ".log"
		}
	}

	lines, err := splitLines(f)
	if err != nil {
		return nil, fmt.Errorf("read lines: %w", err)
	}

	// bucket maps bucket-start (truncated to BucketDuration) → writer + metadata.
	type bucket struct {
		file  *os.File
		result SplitResult
	}
	buckets := map[time.Time]*bucket{}
	var order []time.Time // preserve insertion order for deterministic output

	for _, line := range lines {
		if line.Timestamp.IsZero() {
			continue
		}
		if line.Timestamp.Before(from) || line.Timestamp.After(to) {
			continue
		}

		bucketStart := line.Timestamp.Truncate(opts.BucketDuration)
		b, ok := buckets[bucketStart]
		if !ok {
			bucketEnd := bucketStart.Add(opts.BucketDuration)
			name := fmt.Sprintf("%s_%s%s",
				prefix,
				bucketStart.UTC().Format("20060102T150405Z"),
				ext,
			)
			path := filepath.Join(outDir, name)
			outFile, err := os.Create(path)
			if err != nil {
				return nil, fmt.Errorf("create %s: %w", path, err)
			}
			b = &bucket{
				file: outFile,
				result: SplitResult{
					Path: path,
					From: bucketStart,
					To:   bucketEnd,
				},
			}
			buckets[bucketStart] = b
			order = append(order, bucketStart)
		}

		if _, err := io.WriteString(b.file, line.Raw+"\n"); err != nil {
			return nil, fmt.Errorf("write to %s: %w", b.result.Path, err)
		}
		b.result.LineCount++
	}

	// Close all open files and collect results in insertion order.
	results := make([]SplitResult, 0, len(order))
	for _, key := range order {
		b := buckets[key]
		if err := b.file.Close(); err != nil {
			return nil, fmt.Errorf("close %s: %w", b.result.Path, err)
		}
		results = append(results, b.result)
	}

	return results, nil
}
