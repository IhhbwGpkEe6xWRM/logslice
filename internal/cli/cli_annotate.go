package cli

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/logslice/logslice/internal/slicer"
)

// RunAnnotate is the entry point for the `logslice annotate` subcommand.
func RunAnnotate(args []string) error {
	fs := flag.NewFlagSet("annotate", flag.ContinueOnError)
	input := fs.String("input", "", "input log file (required)")
	fromStr := fs.String("from", "", "start time RFC3339 (required)")
	toStr := fs.String("to", "", "end time RFC3339 (required)")
	output := fs.String("output", "", "output file (default: stdout)")
	lineNumbers := fs.Bool("line-numbers", false, "prefix each line with its line number")
	relTime := fs.Bool("rel-time", false, "prefix each line with ms since first matched line")
	offsets := fs.Bool("offsets", false, "prefix each line with its byte offset")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *input == "" {
		return fmt.Errorf("--input is required")
	}
	if *fromStr == "" {
		return fmt.Errorf("--from is required")
	}
	if *toStr == "" {
		return fmt.Errorf("--to is required")
	}

	from, err := time.Parse(time.RFC3339, *fromStr)
	if err != nil {
		return fmt.Errorf("invalid --from: %w", err)
	}
	to, err := time.Parse(time.RFC3339, *toStr)
	if err != nil {
		return fmt.Errorf("invalid --to: %w", err)
	}
	if from.After(to) {
		return fmt.Errorf("--from must not be after --to")
	}

	opts := slicer.AnnotateOptions{
		AddLineNumbers:  *lineNumbers,
		AddRelativeTime: *relTime,
		AddOffset:       *offsets,
	}

	if *output == "" {
		_, err = slicer.AnnotateFile(*input, from, to, opts, os.Stdout)
		return err
	}
	_, err = slicer.AnnotateFileToFile(*input, *output, from, to, opts)
	return err
}
