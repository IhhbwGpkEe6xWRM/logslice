package cli

import (
	"fmt"
	"io"
	"os"
	"time"

	"logslice/internal/slicer"
)

// AutoDetectArgs holds arguments for auto-detecting format slice.
type AutoDetectArgs struct {
	Input      string
	Output     string
	From       time.Time
	To         time.Time
	SampleSize int
}

// RunAutoDetect opens the input file, auto-detects the timestamp format,
// slices the log to the given time range, and writes output.
func RunAutoDetect(args AutoDetectArgs, stderr io.Writer) error {
	f, err := os.Open(args.Input)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer f.Close()

	var w io.Writer
	if args.Output == "" || args.Output == "-" {
		w = os.Stdout
	} else {
		out, err := os.Create(args.Output)
		if err != nil {
			return fmt.Errorf("create output: %w", err)
		}
		defer out.Close()
		w = out
	}

	opts := slicer.AutoSliceOptions{
		From:       args.From,
		To:         args.To,
		SampleSize: args.SampleSize,
	}

	if err := slicer.AutoSlice(f, w, opts); err != nil {
		return fmt.Errorf("autodetect slice: %w", err)
	}

	if stderr != nil {
		fmt.Fprintf(stderr, "format auto-detected, slice complete\n")
	}
	return nil
}
