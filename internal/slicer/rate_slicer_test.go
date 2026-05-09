package slicer

import (
	"strings"
	"testing"
	"time"
)

func buildRateLog(lines []string) *strings.Reader {
	return strings.NewReader(strings.Join(lines, "\n") + "\n")
}

func mustRateTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSliceRateGroupsByMinute(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:05Z INFO request received",
		"2024-01-01T10:00:45Z INFO request received",
		"2024-01-01T10:01:10Z INFO request received",
		"2024-01-01T10:02:00Z INFO request received",
	}
	from := mustRateTime("2024-01-01T10:00:00Z")
	to := mustRateTime("2024-01-01T10:05:00Z")

	windows, err := SliceRate(buildRateLog(lines), from, to, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(windows))
	}
	if windows[0].Count != 2 {
		t.Errorf("bucket 0: expected count 2, got %d", windows[0].Count)
	}
	if windows[1].Count != 1 {
		t.Errorf("bucket 1: expected count 1, got %d", windows[1].Count)
	}
	if windows[2].Count != 1 {
		t.Errorf("bucket 2: expected count 1, got %d", windows[2].Count)
	}
}

func TestSliceRateExcludesOutOfRange(t *testing.T) {
	lines := []string{
		"2024-01-01T09:59:00Z INFO before range",
		"2024-01-01T10:00:30Z INFO in range",
		"2024-01-01T10:05:01Z INFO after range",
	}
	from := mustRateTime("2024-01-01T10:00:00Z")
	to := mustRateTime("2024-01-01T10:05:00Z")

	windows, err := SliceRate(buildRateLog(lines), from, to, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(windows))
	}
	if windows[0].Count != 1 {
		t.Errorf("expected count 1, got %d", windows[0].Count)
	}
}

func TestSliceRateEmptyInput(t *testing.T) {
	from := mustRateTime("2024-01-01T10:00:00Z")
	to := mustRateTime("2024-01-01T10:05:00Z")

	windows, err := SliceRate(strings.NewReader(""), from, to, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 0 {
		t.Errorf("expected 0 buckets, got %d", len(windows))
	}
}

func TestSliceRateDefaultWindow(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:10Z INFO a",
		"2024-01-01T10:00:50Z INFO b",
	}
	from := mustRateTime("2024-01-01T10:00:00Z")
	to := mustRateTime("2024-01-01T10:05:00Z")

	// window=0 should default to 1 minute
	windows, err := SliceRate(buildRateLog(lines), from, to, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 1 {
		t.Fatalf("expected 1 bucket with default window, got %d", len(windows))
	}
	if windows[0].Count != 2 {
		t.Errorf("expected count 2, got %d", windows[0].Count)
	}
}
