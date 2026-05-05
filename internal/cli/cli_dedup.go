package cli

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/user/logslice/internal/slicer"
)

// RunDedup is the entry-point for the dedup sub-command.
// Usage: logslice dedup -input FILE -from TIME -to TIME [-window N]
func RunDedup(args []string) error {
	fs := flag.NewFlagSet("dedup", flag.ContinueOnError)
	inputFlag := fs.String("input", "", "path to input log file (required)")
	fromFlag := fs.String("from", "", "start time RFC3339 (required)")
	toFlag := fs.String("to", "", "end time RFC3339 (required)")
	windowFlag := fs.Int("window", 1000, "rolling dedup window size")
	outputFlag := fs.String("output", "", "output file path (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *inputFlag == "" {
		return fmt.Errorf("flag -input is required")
	}
	if *fromFlag == "" {
		return fmt.Errorf("flag -from is required")
	}
	if *toFlag == "" {
		return fmt.Errorf("flag -to is required")
	}

	from, err := time.Parse(time.RFC3339, *fromFlag)
	if err != nil {
		return fmt.Errorf("invalid -from: %w", err)
	}
	to, err := time.Parse(time.RFC3339, *toFlag)
	if err != nil {
		return fmt.Errorf("invalid -to: %w", err)
	}
	if from.After(to) {
		return fmt.Errorf("-from must not be after -to")
	}

	f, err := os.Open(*inputFlag)
	if err != nil {
		return fmt.Errorf("cannot open input: %w", err)
	}
	defer f.Close()

	w := os.Stdout
	if *outputFlag != "" {
		out, err := os.Create(*outputFlag)
		if err != nil {
			return fmt.Errorf("cannot create output file: %w", err)
		}
		defer out.Close()
		w = out
	}

	opts := slicer.DedupOptions{WindowSize: *windowFlag}
	n, err := slicer.SliceWithDedup(f, w, from, to, opts)
	if err != nil {
		return fmt.Errorf("dedup slice failed: %w", err)
	}
	fmt.Fprintf(os.Stderr, "wrote %d unique lines\n", n)
	return nil
}
