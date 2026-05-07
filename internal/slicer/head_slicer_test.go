package slicer

import (
	"strings"
	"testing"
	"time"
)

func buildHeadLog() string {
	return strings.Join([]string{
		"2024-01-01T10:00:00Z INFO line-one",
		"2024-01-01T10:00:10Z INFO line-two",
		"2024-01-01T10:00:20Z INFO line-three",
		"2024-01-01T10:00:30Z INFO line-four",
		"2024-01-01T10:00:40Z INFO line-five",
	}, "\n")
}

func TestSliceWithHeadLinesOnly(t *testing.T) {
	from := mustTime("2024-01-01T10:00:00Z")
	to := mustTime("2024-01-01T11:00:00Z")

	lines, err := SliceWithHead(strings.NewReader(buildHeadLog()), from, to, HeadOptions{Lines: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0].Raw != "2024-01-01T10:00:00Z INFO line-one" {
		t.Errorf("unexpected first line: %s", lines[0].Raw)
	}
}

func TestSliceWithHeadDurationOnly(t *testing.T) {
	from := mustTime("2024-01-01T10:00:00Z")
	to := mustTime("2024-01-01T11:00:00Z")

	lines, err := SliceWithHead(strings.NewReader(buildHeadLog()), from, to, HeadOptions{Duration: 25 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// lines at 0s, 10s, 20s fit within 25s window; 30s does not
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestSliceWithHeadCombined(t *testing.T) {
	from := mustTime("2024-01-01T10:00:00Z")
	to := mustTime("2024-01-01T11:00:00Z")

	// Lines=2 is more restrictive than duration=60s
	lines, err := SliceWithHead(strings.NewReader(buildHeadLog()), from, to, HeadOptions{Lines: 2, Duration: 60 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestSliceWithHeadEmptyInput(t *testing.T) {
	from := mustTime("2024-01-01T10:00:00Z")
	to := mustTime("2024-01-01T11:00:00Z")

	lines, err := SliceWithHead(strings.NewReader(""), from, to, HeadOptions{Lines: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
}

func TestSliceWithHeadOutOfRange(t *testing.T) {
	from := mustTime("2024-01-01T12:00:00Z")
	to := mustTime("2024-01-01T13:00:00Z")

	lines, err := SliceWithHead(strings.NewReader(buildHeadLog()), from, to, HeadOptions{Lines: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
}
