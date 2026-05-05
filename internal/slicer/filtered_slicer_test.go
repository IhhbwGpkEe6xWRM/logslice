package slicer

import (
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func mustFilter(t *testing.T, level, pattern string) *parser.LineFilter {
	t.Helper()
	f, err := parser.NewLineFilter(level, pattern)
	if err != nil {
		t.Fatalf("NewLineFilter(%q, %q): %v", level, pattern, err)
	}
	return f
}

func buildFilteredLog() string {
	return strings.Join([]string{
		"2024-03-01T10:00:00Z INFO  startup complete",
		"2024-03-01T10:01:00Z ERROR disk read failed",
		"2024-03-01T10:02:00Z WARN  high memory usage",
		"2024-03-01T10:03:00Z ERROR network timeout",
		"2024-03-01T10:04:00Z INFO  request handled",
	}, "\n") + "\n"
}

func TestSliceWithFilterLevelOnly(t *testing.T) {
	from := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 10, 5, 0, 0, time.UTC)

	var out strings.Builder
	res, err := SliceWithFilter(strings.NewReader(buildFilteredLog()), &out, FilteredSliceOptions{
		From:   from,
		To:     to,
		Filter: mustFilter(t, "ERROR", ""),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LinesMatched != 5 {
		t.Errorf("LinesMatched = %d, want 5", res.LinesMatched)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("output lines = %d, want 2; got:\n%s", len(lines), out.String())
	}
}

func TestSliceWithFilterPatternOnly(t *testing.T) {
	from := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 10, 5, 0, 0, time.UTC)

	var out strings.Builder
	_, err := SliceWithFilter(strings.NewReader(buildFilteredLog()), &out, FilteredSliceOptions{
		From:   from,
		To:     to,
		Filter: mustFilter(t, "", `timeout|memory`),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("output lines = %d, want 2; got:\n%s", len(lines), out.String())
	}
}

func TestSliceWithFilterNilFilter(t *testing.T) {
	from := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 10, 5, 0, 0, time.UTC)

	var out strings.Builder
	res, err := SliceWithFilter(strings.NewReader(buildFilteredLog()), &out, FilteredSliceOptions{
		From:   from,
		To:     to,
		Filter: nil,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LinesFiltered != 0 {
		t.Errorf("LinesFiltered = %d, want 0", res.LinesFiltered)
	}
	if res.LinesMatched != 5 {
		t.Errorf("LinesMatched = %d, want 5", res.LinesMatched)
	}
}
