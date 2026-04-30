package parser

import (
	"bufio"
	"io"
	"strings"
)

// Format represents a detected log timestamp format.
type Format int

const (
	FormatUnknown Format = iota
	FormatRFC3339
	FormatSpaceSeparated
	FormatMilliseconds
)

// String returns a human-readable name for the format.
func (f Format) String() string {
	switch f {
	case FormatRFC3339:
		return "RFC3339"
	case FormatSpaceSeparated:
		return "space-separated"
	case FormatMilliseconds:
		return "milliseconds"
	default:
		return "unknown"
	}
}

// DetectFormat reads up to maxProbeLines lines from r and returns the most
// likely timestamp format found. It does not consume r beyond the probe
// window — callers should pass a fresh reader or use a TeeReader.
func DetectFormat(r io.Reader, maxProbeLines int) Format {
	if maxProbeLines <= 0 {
		maxProbeLines = 20
	}

	counts := make(map[Format]int)

	scanner := bufio.NewScanner(r)
	lines := 0
	for scanner.Scan() && lines < maxProbeLines {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		if f := probeFormat(line); f != FormatUnknown {
			counts[f]++
		}
		lines++
	}

	var best Format
	var bestCount int
	for f, c := range counts {
		if c > bestCount {
			best = f
			bestCount = c
		}
	}
	return best
}

// probeFormat attempts to identify the timestamp format of a single log line.
func probeFormat(line string) Format {
	if len(line) < 10 {
		return FormatUnknown
	}
	// RFC3339: 2006-01-02T15:04:05
	if len(line) >= 19 && line[4] == '-' && line[7] == '-' && line[10] == 'T' {
		return FormatRFC3339
	}
	// Space-separated: 2006-01-02 15:04:05
	if len(line) >= 19 && line[4] == '-' && line[7] == '-' && line[10] == ' ' && line[13] == ':' {
		return FormatSpaceSeparated
	}
	// Milliseconds prefix: 13-digit unix ms timestamp followed by space
	if len(line) >= 14 && line[13] == ' ' {
		allDigits := true
		for _, c := range line[:13] {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			return FormatMilliseconds
		}
	}
	return FormatUnknown
}
