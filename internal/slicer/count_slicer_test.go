package slicer

import (
	"strings"
	"testing"
	"time"
)

func buildCountLog() string {
	return strings.Join([]string{
		`2024-01-15T10:00:00Z INFO  user logged in`,
		`2024-01-15T10:01:00Z ERROR disk full`,
		`2024-01-15T10:01:30Z WARN  high memory`,
		`2024-01-15T10:02:00Z INFO  backup started`,
		`2024-01-15T10:03:00Z ERROR connection refused`,
		`2024-01-15T11:00:00Z INFO  outside range`,
		`not a log line`,
	}, "\n")
}

func mustCountTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSliceCountMatchedLines(t *testing.T) {
	r := strings.NewReader(buildCountLog())
	from := mustCountTime("2024-01-15T10:00:00Z")
	to := mustCountTime("2024-01-15T10:03:00Z")

	result, err := SliceCount(r, from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Matched != 5 {
		t.Errorf("expected 5 matched, got %d", result.Matched)
	}
	if result.Total != 6 {
		t.Errorf("expected 6 total (excluding blank/non-log), got %d", result.Total)
	}
}

func TestSliceCountByLevel(t *testing.T) {
	r := strings.NewReader(buildCountLog())
	from := mustCountTime("2024-01-15T10:00:00Z")
	to := mustCountTime("2024-01-15T10:03:00Z")

	result, err := SliceCount(r, from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ByLevel["ERROR"] != 2 {
		t.Errorf("expected 2 ERROR lines, got %d", result.ByLevel["ERROR"])
	}
	if result.ByLevel["INFO"] != 2 {
		t.Errorf("expected 2 INFO lines, got %d", result.ByLevel["INFO"])
	}
	if result.ByLevel["WARN"] != 1 {
		t.Errorf("expected 1 WARN line, got %d", result.ByLevel["WARN"])
	}
}

func TestSliceCountByMinute(t *testing.T) {
	r := strings.NewReader(buildCountLog())
	from := mustCountTime("2024-01-15T10:00:00Z")
	to := mustCountTime("2024-01-15T10:03:00Z")

	result, err := SliceCount(r, from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ByMinute["2024-01-15T10:01"] != 2 {
		t.Errorf("expected 2 lines at 10:01, got %d", result.ByMinute["2024-01-15T10:01"])
	}
}

func TestSliceCountEmptyInput(t *testing.T) {
	r := strings.NewReader("")
	from := mustCountTime("2024-01-15T10:00:00Z")
	to := mustCountTime("2024-01-15T11:00:00Z")

	result, err := SliceCount(r, from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 0 || result.Matched != 0 {
		t.Errorf("expected zero counts, got total=%d matched=%d", result.Total, result.Matched)
	}
}
