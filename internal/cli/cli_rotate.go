package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/logslice/logslice/internal/slicer"
)

// RunRotate is the entry point for the `logslice rotate` subcommand.
// It slices across a directory of rotated log files and writes matching
// lines to stdout or an output file.
func RunRotate(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("rotate", flag.ContinueOnError)
	fs.SetOutput(stdout)

	dir := fs.String("dir", "", "directory containing rotated log files (required)")
	pattern := fs.String("pattern", "*.log*", "glob pattern to match log files")
	fromStr := fs.String("from", "", "start of time range, RFC3339 (required)")
	toStr := fs.String("to", "", "end of time range, RFC3339 (optional)")
	output := fs.String("output", "", "output file path (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *dir == "" {
		return fmt.Errorf("rotate: -dir is required")
	}
	if *fromStr == "" {
		return fmt.Errorf("rotate: -from is required")
	}

	from, err := time.Parse(time.RFC3339, *fromStr)
	if err != nil {
		return fmt.Errorf("rotate: invalid -from %q: %w", *fromStr, err)
	}

	var to time.Time
	if *toStr != "" {
		to, err = time.Parse(time.RFC3339, *toStr)
		if err != nil {
			return fmt.Errorf("rotate: invalid -to %q: %w", *toStr, err)
		}
		if from.After(to) {
			return fmt.Errorf("rotate: -from %s is after -to %s", from, to)
		}
	}

	var w io.Writer = stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("rotate: create output %q: %w", *output, err)
		}
		defer f.Close()
		w = f
	}

	opts := slicer.RotateSliceOptions{
		Dir:     *dir,
		Pattern: *pattern,
		From:    from,
		To:      to,
	}

	if err := slicer.RotateSlice(opts, w); err != nil {
		return fmt.Errorf("rotate: %w", err)
	}
	return nil
}
