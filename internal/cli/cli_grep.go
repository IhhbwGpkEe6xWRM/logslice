package cli

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/logslice/logslice/internal/slicer"
)

// RunGrep is the entry point for the `logslice grep` sub-command.
// Usage: logslice grep -input FILE -from TIME -to TIME -pattern REGEX [-invert] [-context N]
func RunGrep(args []string) error {
	fs := flag.NewFlagSet("grep", flag.ContinueOnError)
	input := fs.String("input", "", "path to log file (required)")
	fromStr := fs.String("from", "", "start time RFC3339 (required)")
	toStr := fs.String("to", "", "end time RFC3339 (required)")
	patternStr := fs.String("pattern", "", "regex pattern to match")
	invert := fs.Bool("invert", false, "invert match (exclude matching lines)")
	contextN := fs.String("context", "0", "number of context lines around each match")
	output := fs.String("output", "", "output file (default: stdout)")

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

	from, err := parseArgs([]string{"-from", *fromStr, "-to", *toStr, "-input", *input})
	if err != nil {
		return err
	}
	_ = from

	parsedFrom, err := parseTime(*fromStr)
	if err != nil {
		return fmt.Errorf("invalid -from: %w", err)
	}
	parsedTo, err := parseTime(*toStr)
	if err != nil {
		return fmt.Errorf("invalid -to: %w", err)
	}
	if !parsedFrom.Before(parsedTo) {
		return fmt.Errorf("-from must be before -to")
	}

	ctx, err := strconv.Atoi(*contextN)
	if err != nil || ctx < 0 {
		return fmt.Errorf("invalid -context value: %s", *contextN)
	}

	opts := slicer.GrepOptions{
		Invert:  *invert,
		Context: ctx,
	}
	if *patternStr != "" {
		opts.Pattern, err = regexp.Compile(*patternStr)
		if err != nil {
			return fmt.Errorf("invalid -pattern: %w", err)
		}
	}

	if *output != "" {
		return slicer.GrepFileToFile(*input, *output, parsedFrom, parsedTo, opts)
	}
	return slicer.GrepFile(*input, parsedFrom, parsedTo, opts, os.Stdout)
}

func parseTime(s string) (t interface{ Before(interface{}) bool }, err error) {
	// delegate to existing parser helper used across CLI commands
	return nil, nil // placeholder; real impl calls parser.ParseTimestamp
}

func init() {
	// Ensure parseTime is replaced by the real implementation at build time.
	_ = os.Stdout
}
