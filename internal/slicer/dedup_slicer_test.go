package slicer

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func buildDedupLog(entries []string) string {
	return strings.Join(entries, "\n") + "\n"
}

func TestSliceWithDedupRemovesDuplicates(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)

	lines := []string{
		"2024-01-01T00:10:00Z INFO user logged in",
		"2024-01-01T00:20:00Z INFO user logged in",
		"2024-01-01T00:30:00Z ERROR disk full",
		"2024-01-01T00:40:00Z INFO user logged in",
	}
	input := buildDedupLog(lines)

	var buf bytes.Buffer
	n, err := SliceWithDedup(strings.NewReader(input), &buf, from, to, DedupOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 unique lines, got %d", n)
	}
	out := buf.String()
	if !strings.Contains(out, "user logged in") {
		t.Error("expected first occurrence of 'user logged in' in output")
	}
	if !strings.Contains(out, "disk full") {
		t.Error("expected 'disk full' in output")
	}
}

func TestSliceWithDedupWindowEviction(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 2, 0, 0, 0, time.UTC)

	// Window of 2: after 2 unique messages, the first is evicted and a
	// repeat of it should be written again.
	lines := []string{
		"2024-01-01T00:01:00Z INFO msg A",
		"2024-01-01T00:02:00Z INFO msg B",
		"2024-01-01T00:03:00Z INFO msg C", // evicts A
		"2024-01-01T00:04:00Z INFO msg A", // A no longer in window → written
	}
	input := buildDedupLog(lines)

	var buf bytes.Buffer
	n, err := SliceWithDedup(strings.NewReader(input), &buf, from, to, DedupOptions{WindowSize: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 4 {
		t.Errorf("expected 4 written lines with window=2, got %d", n)
	}
}

func TestSliceWithDedupEmptyInput(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)

	var buf bytes.Buffer
	n, err := SliceWithDedup(strings.NewReader(""), &buf, from, to, DedupOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
}

func TestSliceWithDedupOutOfRange(t *testing.T) {
	from := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)

	lines := []string{
		"2024-01-01T00:01:00Z INFO early message",
		"2024-01-01T00:02:00Z INFO another early message",
	}

	var buf bytes.Buffer
	n, err := SliceWithDedup(strings.NewReader(buildDedupLog(lines)), &buf, from, to, DedupOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines out of range, got %d", n)
	}
}
