package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/slicer"
)

// RunSample implements the `logslice sample` sub-command.
//
// Usage: logslice sample <input> <from> <to> <rate> [output]
//
//	<rate>   float in (0, 1] – fraction of in-range lines to keep
//	[output] optional output file; defaults to stdout
func RunSample(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("usage: logslice sample <input> <from> <to> <rate> [output]")
	}

	inputPath := args[0]
	fromStr := args[1]
	toStr := args[2]
	rateStr := args[3]

	from, err := parser.ParseTimestamp(fromStr)
	if err != nil {
		return fmt.Errorf("invalid --from timestamp %q: %w", fromStr, err)
	}

	to, err := parser.ParseTimestamp(toStr)
	if err != nil {
		return fmt.Errorf("invalid --to timestamp %q: %w", toStr, err)
	}

	if from.After(to.Time) {
		return fmt.Errorf("--from must not be after --to")
	}

	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil || rate <= 0 || rate > 1.0 {
		return fmt.Errorf("rate must be a float in (0, 1], got %q", rateStr)
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open %q: %w", inputPath, err)
	}
	defer f.Close()

	var w io.Writer = os.Stdout
	if len(args) >= 5 {
		out, err := os.Create(args[4])
		if err != nil {
			return fmt.Errorf("create output %q: %w", args[4], err)
		}
		defer out.Close()
		w = out
	}

	n, err := slicer.SampleSlice(f, w, from, to, slicer.SampleOptions{Rate: rate})
	if err != nil {
		return fmt.Errorf("sample: %w", err)
	}

	fmt.Fprintf(os.Stderr, "sampled %d lines (rate=%.2f)\n", n, rate)
	return nil
}
