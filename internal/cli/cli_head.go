package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/logslice/logslice/internal/slicer"
)

// RunHead is the entry point for the `logslice head` sub-command.
// Usage: logslice head -input <file> -from <ts> -to <ts> [-lines N] [-duration D]
func RunHead(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("head", flag.ContinueOnError)

	input := fs.String("input", "", "path to log file (required)")
	fromStr := fs.String("from", "", "start timestamp RFC3339 (required)")
	toStr := fs.String("to", "", "end timestamp RFC3339 (required)")
	linesStr := fs.String("lines", "0", "maximum number of lines to return")
	durStr := fs.String("duration", "", "maximum duration window from first matched line (e.g. 30s, 5m)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *input == "" {
		return fmt.Errorf("flag -input is required")
	}
	if *fromStr == "" {
		return fmt.Errorf("flag -from is required")
	}
	if *toStr == "" {
		return fmt.Errorf("flag -to is required")
	}

	from, err := time.Parse(time.RFC3339, *fromStr)
	if err != nil {
		return fmt.Errorf("invalid -from: %w", err)
	}
	to, err := time.Parse(time.RFC3339, *toStr)
	if err != nil {
		return fmt.Errorf("invalid -to: %w", err)
	}
	if from.After(to) {
		return fmt.Errorf("-from must not be after -to")
	}

	n, err := strconv.Atoi(*linesStr)
	if err != nil || n < 0 {
		return fmt.Errorf("invalid -lines: must be a non-negative integer")
	}

	var dur time.Duration
	if *durStr != "" {
		dur, err = time.ParseDuration(*durStr)
		if err != nil {
			return fmt.Errorf("invalid -duration: %w", err)
		}
	}

	opts := slicer.HeadOptions{Lines: n, Duration: dur}

	w := stdout
	if w == nil {
		w = os.Stdout
	}

	return slicer.HeadFileToWriter(*input, from, to, opts, w)
}
