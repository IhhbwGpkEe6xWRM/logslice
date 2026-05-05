package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/user/logslice/internal/slicer"
)

// RunSplit splits a log file into multiple files bucketed by a time window.
// Args: <input> <output-dir> <window-duration>
// Example window-duration: 1h, 30m, 24h
func RunSplit(args []string, stderr io.Writer) error {
	if len(args) < 3 {
		return errors.New("usage: logslice split <input> <output-dir> <window>")
	}

	input := args[0]
	outDir := args[1]
	windowStr := args[2]

	win, err := parseSplitWindow(windowStr)
	if err != nil {
		return fmt.Errorf("invalid window %q: %w", windowStr, err)
	}

	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %w", err)
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("cannot create output dir: %w", err)
	}

	n, err := slicer.SplitFile(input, outDir, win)
	if err != nil {
		return fmt.Errorf("split failed: %w", err)
	}

	fmt.Fprintf(stderr, "split complete: %d file(s) written to %s\n", n, outDir)
	return nil
}

// parseSplitWindow accepts Go duration strings (e.g. "1h") or plain integer
// minutes for convenience (e.g. "60" → 60 minutes).
func parseSplitWindow(s string) (time.Duration, error) {
	if d, err := time.ParseDuration(s); err == nil {
		if d <= 0 {
			return 0, errors.New("window must be positive")
		}
		return d, nil
	}
	// fallback: plain integer treated as minutes
	if mins, err := strconv.Atoi(s); err == nil {
		if mins <= 0 {
			return 0, errors.New("window must be positive")
		}
		return time.Duration(mins) * time.Minute, nil
	}
	return 0, fmt.Errorf("cannot parse %q as duration", s)
}
