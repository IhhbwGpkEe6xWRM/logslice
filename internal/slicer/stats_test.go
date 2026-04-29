package slicer

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestSliceWithStatsMatchedLines(t *testing.T) {
	from := mustTime("2024-01-01T10:00:00Z")
	to := mustTime("2024-01-01T10:05:00Z")
	s := New(from, to)

	input := strings.Join([]string{
		`2024-01-01T09:59:00Z level=info msg="before range"`,
		`2024-01-01T10:01:00Z level=info msg="in range"`,
		`2024-01-01T10:03:00Z level=info msg="also in range"`,
		`2024-01-01T10:06:00Z level=info msg="after range"`,
		`not a log line`,
	}, "\n")

	var out bytes.Buffer
	stats, err := s.SliceWithStats(strings.NewReader(input), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.TotalLines != 5 {
		t.Errorf("TotalLines: got %d, want 5", stats.TotalLines)
	}
	if stats.MatchedLines != 2 {
		t.Errorf("MatchedLines: got %d, want 2", stats.MatchedLines)
	}
	if stats.SkippedLines != 2 {
		t.Errorf("SkippedLines: got %d, want 2", stats.SkippedLines)
	}
	if stats.UnparsedLines != 1 {
		t.Errorf("UnparsedLines: got %d, want 1", stats.UnparsedLines)
	}
	if stats.Duration <= 0 {
		t.Error("Duration should be positive")
	}
	if stats.FirstMatch == nil || !stats.FirstMatch.Equal(mustTime("2024-01-01T10:01:00Z")) {
		t.Errorf("FirstMatch: got %v, want 2024-01-01T10:01:00Z", stats.FirstMatch)
	}
	if stats.LastMatch == nil || !stats.LastMatch.Equal(mustTime("2024-01-01T10:03:00Z")) {
		t.Errorf("LastMatch: got %v, want 2024-01-01T10:03:00Z", stats.LastMatch)
	}
}

func TestStatsMatchRate(t *testing.T) {
	s := &Stats{TotalLines: 10, MatchedLines: 4}
	if got := s.MatchRate(); got != 0.4 {
		t.Errorf("MatchRate: got %f, want 0.4", got)
	}
}

func TestStatsMatchRateZeroTotal(t *testing.T) {
	s := &Stats{}
	if got := s.MatchRate(); got != 0 {
		t.Errorf("MatchRate with zero total: got %f, want 0", got)
	}
}

func TestStatsSummary(t *testing.T) {
	ts := mustTime("2024-01-01T10:00:00Z")
	s := &Stats{
		TotalLines:    100,
		MatchedLines:  20,
		SkippedLines:  75,
		UnparsedLines: 5,
		Duration:      42 * time.Millisecond,
		FirstMatch:    &ts,
		LastMatch:     &ts,
	}
	var buf bytes.Buffer
	s.Summary(&buf)
	out := buf.String()
	for _, want := range []string{"100", "20", "75", "5", "42ms", "2024-01-01T10:00:00Z"} {
		if !strings.Contains(out, want) {
			t.Errorf("Summary output missing %q\ngot:\n%s", want, out)
		}
	}
}
