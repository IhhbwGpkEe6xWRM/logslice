package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/user/logslice/internal/slicer"
)

// sliceToWriterWithProgress is like sliceToWriter but reports progress to
// stderr when the input is a regular file (size is known).
func sliceToWriterWithProgress(args parsedArgs, dst io.Writer) error {
	f, err := os.Open(args.inputPath)
	if err != nil {
		return fmt.Errorf("open %q: %w", args.inputPath, err)
	}
	defer f.Close()

	var totalBytes int64
	if info, statErr := f.Stat(); statErr == nil {
		totalBytes = info.Size()
	}

	reporter := slicer.NewProgressReporter(os.Stderr, totalBytes)

	stats, err := slicer.SliceWithProgress(f, dst, args.from, args.to, reporter)
	if err != nil {
		return fmt.Errorf("slice: %w", err)
	}

	fmt.Fprintf(os.Stderr, "done: %d/%d lines matched (%.1f%%)\n",
		stats.MatchedLines, stats.TotalLines, stats.MatchRate()*100)
	return nil
}
