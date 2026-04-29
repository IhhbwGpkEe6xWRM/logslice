package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/slicer"
)

// Config holds parsed CLI arguments.
type Config struct {
	From   time.Time
	To     time.Time
	Input  string
	Output string
}

// Run parses args and executes the slice operation.
func Run(args []string) error {
	cfg, err := parseArgs(args)
	if err != nil {
		return err
	}

	if cfg.Output != "" {
		return slicer.SliceFileToFile(cfg.Input, cfg.Output, cfg.From, cfg.To)
	}

	return sliceToWriter(cfg.Input, cfg.From, cfg.To, os.Stdout)
}

func parseArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	var fromStr, toStr, input, output string
	fs.StringVar(&fromStr, "from", "", "start timestamp (RFC3339 or common log format)")
	fs.StringVar(&toStr, "to", "", "end timestamp (RFC3339 or common log format)")
	fs.StringVar(&input, "input", "", "input log file path (required)")
	fs.StringVar(&output, "output", "", "output file path (optional, defaults to stdout)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if input == "" {
		return nil, errors.New("--input is required")
	}
	if fromStr == "" {
		return nil, errors.New("--from is required")
	}
	if toStr == "" {
		return nil, errors.New("--to is required")
	}

	from, err := parser.ParseTimestamp(fromStr)
	if err != nil {
		return nil, fmt.Errorf("invalid --from: %w", err)
	}
	to, err := parser.ParseTimestamp(toStr)
	if err != nil {
		return nil, fmt.Errorf("invalid --to: %w", err)
	}
	if !from.Before(to) {
		return nil, errors.New("--from must be before --to")
	}

	return &Config{From: from, To: to, Input: input, Output: output}, nil
}

func sliceToWriter(inputPath string, from, to time.Time, w io.Writer) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("opening input: %w", err)
	}
	defer f.Close()

	s := slicer.New(from, to)
	return s.Slice(f, w)
}
