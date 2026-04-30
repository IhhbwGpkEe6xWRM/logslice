package parser_test

import (
	"strings"
	"testing"

	"logslice/internal/parser"
)

func TestSampleReaderRFC3339(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-15T10:00:00Z INFO starting server",
		"2024-01-15T10:00:01Z DEBUG listening on :8080",
		"2024-01-15T10:00:02Z INFO ready",
	}, "\n")

	result := parser.SampleReader(strings.NewReader(input), 20)
	if !result.Detected {
		t.Fatalf("expected format to be detected, got Detected=false")
	}
	if result.Format != "rfc3339" {
		t.Errorf("expected rfc3339, got %q", result.Format)
	}
	if result.SampleSize != 3 {
		t.Errorf("expected SampleSize=3, got %d", result.SampleSize)
	}
}

func TestSampleReaderSpaceSeparated(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-15 10:00:00 INFO starting",
		"2024-01-15 10:00:01 DEBUG ok",
	}, "\n")

	result := parser.SampleReader(strings.NewReader(input), 20)
	if !result.Detected {
		t.Fatalf("expected format to be detected")
	}
	if result.Format != "space" {
		t.Errorf("expected space, got %q", result.Format)
	}
}

func TestSampleReaderEmpty(t *testing.T) {
	result := parser.SampleReader(strings.NewReader(""), 20)
	if result.Detected {
		t.Error("expected Detected=false for empty input")
	}
	if result.SampleSize != 0 {
		t.Errorf("expected SampleSize=0, got %d", result.SampleSize)
	}
}

func TestSampleReaderUnknown(t *testing.T) {
	input := "no timestamp here\njust plain text\nnothing useful"
	result := parser.SampleReader(strings.NewReader(input), 20)
	if result.Detected {
		t.Error("expected Detected=false for lines with no timestamps")
	}
}

func TestSampleReaderRespectsMaxLines(t *testing.T) {
	lines := make([]string, 50)
	for i := range lines {
		lines[i] = "2024-01-15T10:00:00Z INFO line"
	}
	input := strings.Join(lines, "\n")
	result := parser.SampleReader(strings.NewReader(input), 10)
	if result.SampleSize > 10 {
		t.Errorf("expected SampleSize<=10, got %d", result.SampleSize)
	}
}

func TestSampleLinesDetected(t *testing.T) {
	lines := []string{
		"2024-01-15T10:00:00Z INFO a",
		"2024-01-15T10:00:01Z INFO b",
	}
	result := parser.SampleLines(lines)
	if !result.Detected {
		t.Error("expected Detected=true")
	}
	if result.SampleSize != 2 {
		t.Errorf("expected SampleSize=2, got %d", result.SampleSize)
	}
}

func TestSampleLinesEmpty(t *testing.T) {
	result := parser.SampleLines(nil)
	if result.Detected {
		t.Error("expected Detected=false for nil input")
	}
}
