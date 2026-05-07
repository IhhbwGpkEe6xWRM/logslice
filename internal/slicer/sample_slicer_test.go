package slicer

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

func buildSampleLog(from, to time.Time, total int) string {
	var sb strings.Builder
	step := to.Sub(from) / time.Duration(total)
	for i := 0; i < total; i++ {
		t := from.Add(time.Duration(i) * step)
		sb.WriteString(t.UTC().Format(time.RFC3339) + " sample line " + string(rune('A'+i%26)) + "\n")
	}
	return sb.String()
}

func mustTS(s string) parser.LogTimestamp {
	return parser.MustParseTimestamp(s)
}

func TestSampleSliceRate100(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	input := buildSampleLog(from, to, 20)

	r := strings.NewReader(input)
	var w bytes.Buffer

	n, err := SampleSlice(r, &w, mustTS("2024-01-01T00:00:00Z"), mustTS("2024-01-01T01:00:00Z"),
		SampleOptions{Rate: 1.0, Seed: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 20 {
		t.Errorf("expected 20 lines, got %d", n)
	}
}

func TestSampleSliceRateHalf(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	input := buildSampleLog(from, to, 1000)

	r := strings.NewReader(input)
	var w bytes.Buffer

	n, err := SampleSlice(r, &w, mustTS("2024-01-01T00:00:00Z"), mustTS("2024-01-01T01:00:00Z"),
		SampleOptions{Rate: 0.5, Seed: 7})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// With seed 7 and 1000 lines at 50% we expect roughly 400-600.
	if n < 400 || n > 600 {
		t.Errorf("expected ~500 lines, got %d", n)
	}
}

func TestSampleSliceEmptyInput(t *testing.T) {
	r := strings.NewReader("")
	var w bytes.Buffer

	n, err := SampleSlice(r, &w, mustTS("2024-01-01T00:00:00Z"), mustTS("2024-01-01T01:00:00Z"),
		SampleOptions{Rate: 1.0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
}

func TestSampleSliceOutOfRange(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	input := buildSampleLog(from, to, 10)

	r := strings.NewReader(input)
	var w bytes.Buffer

	// Query a window that doesn't overlap with the log.
	n, err := SampleSlice(r, &w, mustTS("2024-01-02T00:00:00Z"), mustTS("2024-01-02T01:00:00Z"),
		SampleOptions{Rate: 1.0, Seed: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
}

func TestSampleSliceInvalidRateDefaultsToFull(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	input := buildSampleLog(from, to, 10)

	r := strings.NewReader(input)
	var w bytes.Buffer

	// Rate of 0 should default to 1.0 (keep all).
	n, err := SampleSlice(r, &w, mustTS("2024-01-01T00:00:00Z"), mustTS("2024-01-01T01:00:00Z"),
		SampleOptions{Rate: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 10 {
		t.Errorf("expected 10 lines, got %d", n)
	}
}
