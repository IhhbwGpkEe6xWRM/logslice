package slicer

import (
	"strings"
	"testing"
	"time"
)

func buildContextLog() string {
	return strings.Join([]string{
		"2024-01-01T10:00:00Z INFO  before-before",
		"2024-01-01T10:01:00Z INFO  before",
		"2024-01-01T10:02:00Z INFO  match-start",
		"2024-01-01T10:03:00Z INFO  match-middle",
		"2024-01-01T10:04:00Z INFO  match-end",
		"2024-01-01T10:05:00Z INFO  after",
		"2024-01-01T10:06:00Z INFO  after-after",
	}, "\n")
}

func mustContextTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSliceWithContextNoContext(t *testing.T) {
	r := strings.NewReader(buildContextLog())
	from := mustContextTime("2024-01-01T10:02:00Z")
	to := mustContextTime("2024-01-01T10:04:00Z")

	lines, err := SliceWithContext(r, from, to, ContextOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0].Raw != "2024-01-01T10:02:00Z INFO  match-start" {
		t.Errorf("unexpected first line: %q", lines[0].Raw)
	}
}

func TestSliceWithContextBefore(t *testing.T) {
	r := strings.NewReader(buildContextLog())
	from := mustContextTime("2024-01-01T10:02:00Z")
	to := mustContextTime("2024-01-01T10:02:00Z")

	lines, err := SliceWithContext(r, from, to, ContextOptions{Before: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (2 before + match), got %d", len(lines))
	}
	if lines[0].Raw != "2024-01-01T10:00:00Z INFO  before-before" {
		t.Errorf("unexpected first line: %q", lines[0].Raw)
	}
}

func TestSliceWithContextAfter(t *testing.T) {
	r := strings.NewReader(buildContextLog())
	from := mustContextTime("2024-01-01T10:04:00Z")
	to := mustContextTime("2024-01-01T10:04:00Z")

	lines, err := SliceWithContext(r, from, to, ContextOptions{After: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (match + 2 after), got %d", len(lines))
	}
	if lines[2].Raw != "2024-01-01T10:06:00Z INFO  after-after" {
		t.Errorf("unexpected last line: %q", lines[2].Raw)
	}
}

func TestSliceWithContextEmptyInput(t *testing.T) {
	r := strings.NewReader("")
	from := mustContextTime("2024-01-01T10:00:00Z")
	to := mustContextTime("2024-01-01T11:00:00Z")

	lines, err := SliceWithContext(r, from, to, ContextOptions{Before: 2, After: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
}

func TestSliceWithContextOverlappingWindows(t *testing.T) {
	r := strings.NewReader(buildContextLog())
	from := mustContextTime("2024-01-01T10:02:00Z")
	to := mustContextTime("2024-01-01T10:04:00Z")

	lines, err := SliceWithContext(r, from, to, ContextOptions{Before: 1, After: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// before(1) + 3 matches + after(1) = 5, no duplicates
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(lines))
	}
}
