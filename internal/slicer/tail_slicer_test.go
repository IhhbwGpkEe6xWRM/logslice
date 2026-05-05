package slicer

import (
	"strings"
	"testing"
	"time"
)

func buildTailLog(entries []string) *strings.Reader {
	return strings.NewReader(strings.Join(entries, "\n"))
}

func TestSliceWithTailLinesOnly(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z info line1",
		"2024-01-01T10:01:00Z info line2",
		"2024-01-01T10:02:00Z info line3",
		"2024-01-01T10:03:00Z info line4",
	}
	r := buildTailLog(lines)
	res, err := SliceWithTail(r, TailOptions{Lines: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(res.Lines))
	}
	if res.Lines[0] != lines[2] || res.Lines[1] != lines[3] {
		t.Errorf("unexpected lines: %v", res.Lines)
	}
	if res.Stats.TotalLines != 4 || res.Stats.MatchedLines != 2 {
		t.Errorf("unexpected stats: %+v", res.Stats)
	}
}

func TestSliceWithTailDurationOnly(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z info old",
		"2024-01-01T10:04:00Z info recent",
		"2024-01-01T10:05:00Z info latest",
	}
	r := buildTailLog(lines)
	res, err := SliceWithTail(r, TailOptions{Duration: 2 * time.Minute})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(res.Lines))
	}
}

func TestSliceWithTailCombined(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z info a",
		"2024-01-01T10:04:00Z info b",
		"2024-01-01T10:05:00Z info c",
		"2024-01-01T10:05:30Z info d",
	}
	r := buildTailLog(lines)
	// Duration keeps last 2 min (b,c,d), Lines cap at 2 => (c,d)
	res, err := SliceWithTail(r, TailOptions{Lines: 2, Duration: 2 * time.Minute})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(res.Lines))
	}
	if res.Lines[1] != lines[3] {
		t.Errorf("last line should be d, got %q", res.Lines[1])
	}
}

func TestSliceWithTailEmptyInput(t *testing.T) {
	r := strings.NewReader("")
	res, err := SliceWithTail(r, TailOptions{Lines: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Lines) != 0 {
		t.Errorf("expected no lines, got %d", len(res.Lines))
	}
}

func TestSliceWithTailNoOptions(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z info x",
		"2024-01-01T10:01:00Z info y",
	}
	r := buildTailLog(lines)
	res, err := SliceWithTail(r, TailOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Lines) != 2 {
		t.Errorf("expected all 2 lines, got %d", len(res.Lines))
	}
}
