package slicer

import (
	"bufio"
	"io"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// UniqueField specifies which part of a log line to deduplicate on.
type UniqueField int

const (
	UniqueByMessage UniqueField = iota
	UniqueByLevel
	UniqueByFull
)

// UniqueOptions controls the behaviour of SliceUnique.
type UniqueOptions struct {
	From   time.Time
	To     time.Time
	Field  UniqueField
	Format string // optional forced format; empty = auto
}

// SliceUnique reads lines from r, keeps only the first occurrence of each
// unique key (determined by Field) within the time range [From, To], and
// writes matching lines to w.
func SliceUnique(r io.Reader, w io.Writer, opts UniqueOptions) (int, error) {
	seen := make(map[string]struct{})
	written := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		raw := scanner.Text()
		if raw == "" {
			continue
		}

		line := parser.ParseLine(raw)
		if line.Timestamp.IsZero() {
			continue
		}
		if !line.InRange(opts.From, opts.To) {
			// Past the end of range — we can stop early.
			if !line.Timestamp.Before(opts.To.Add(time.Second)) && !opts.To.IsZero() {
				break
			}
			continue
		}

		key := uniqueKey(line, opts.Field)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		if _, err := io.WriteString(w, raw+"\n"); err != nil {
			return written, err
		}
		written++
	}

	return written, scanner.Err()
}

func uniqueKey(line parser.LogLine, field UniqueField) string {
	switch field {
	case UniqueByLevel:
		return line.Level
	case UniqueByFull:
		return line.Raw
	default: // UniqueByMessage
		return line.Message
	}
}
