package slicer

import (
	"bufio"
	"io"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// HeadOptions controls how many lines or how much duration to take from the start.
type HeadOptions struct {
	Lines    int
	Duration time.Duration
}

// SliceWithHead reads lines from r within [from, to] and returns at most the
// first N lines (by count and/or duration window from the first matched line).
func SliceWithHead(r io.Reader, from, to time.Time, opts HeadOptions) ([]parser.LogLine, error) {
	scanner := bufio.NewScanner(r)

	var (
		results  []parser.LogLine
		firstTS  time.Time
		count    int
	)

	for scanner.Scan() {
		text := scanner.Text()
		line, err := parser.ParseLine(text)
		if err != nil || line.Timestamp.IsZero() {
			continue
		}
		if !line.InRange(from, to) {
			if !line.Timestamp.Before(from) {
				break
			}
			continue
		}

		if firstTS.IsZero() {
			firstTS = line.Timestamp
		}

		if opts.Duration > 0 && line.Timestamp.Sub(firstTS) > opts.Duration {
			break
		}

		results = append(results, line)
		count++

		if opts.Lines > 0 && count >= opts.Lines {
			break
		}
	}

	return results, scanner.Err()
}

// filterHead applies head options to an already-collected slice of lines.
func filterHead(lines []parser.LogLine, opts HeadOptions) []parser.LogLine {
	if len(lines) == 0 {
		return lines
	}
	firstTS := lines[0].Timestamp
	var out []parser.LogLine
	for _, l := range lines {
		if opts.Duration > 0 && l.Timestamp.Sub(firstTS) > opts.Duration {
			break
		}
		out = append(out, l)
		if opts.Lines > 0 && len(out) >= opts.Lines {
			break
		}
	}
	return out
}
