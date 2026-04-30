package parser

import (
	"bufio"
	"io"
)

// SampleResult holds the outcome of sampling a log file for format detection.
type SampleResult struct {
	Format     string
	SampleSize int
	Detected   bool
}

// SampleReader reads up to maxLines lines from r and attempts to detect
// the timestamp format used in the log file.
func SampleReader(r io.Reader, maxLines int) SampleResult {
	if maxLines <= 0 {
		maxLines = 20
	}

	scanner := bufio.NewScanner(r)
	lines := make([]string, 0, maxLines)
	for scanner.Scan() && len(lines) < maxLines {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		return SampleResult{}
	}

	fmt := DetectFormat(lines)
	return SampleResult{
		Format:   fmt,
		SampleSize: len(lines),
		Detected: fmt != "unknown" && fmt != "",
	}
}

// SampleLines attempts to detect the timestamp format from a slice of log lines.
func SampleLines(lines []string) SampleResult {
	if len(lines) == 0 {
		return SampleResult{}
	}
	fmt := DetectFormat(lines)
	return SampleResult{
		Format:   fmt,
		SampleSize: len(lines),
		Detected: fmt != "unknown" && fmt != "",
	}
}
