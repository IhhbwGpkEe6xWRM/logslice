package slicer

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func buildProgressLog() string {
	lines := []string{
		`2024-01-10T10:00:00Z INFO start`,
		`2024-01-10T10:05:00Z DEBUG tick`,
		`2024-01-10T10:10:00Z INFO middle`,
		`2024-01-10T10:15:00Z WARN warn`,
		`2024-01-10T10:20:00Z INFO end`,
	}
	return strings.Join(lines, "\n") + "\n"
}

func TestSliceWithProgressMatchesRange(t *testing.T) {
	from := mustTime("2024-01-10T10:05:00Z")
	to := mustTime("2024-01-10T10:15:00Z")

	src := strings.NewReader(buildProgressLog())
	var dst bytes.Buffer
	var progBuf bytes.Buffer

	reporter := NewProgressReporter(&progBuf, int64(len(buildProgressLog())))
	stats, err := SliceWithProgress(src, &dst, from, to, reporter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := dst.String()
	if !strings.Contains(got, "10:05:00") {
		t.Errorf("expected 10:05 line in output")
	}
	if strings.Contains(got, "10:00:00") {
		t.Errorf("expected 10:00 line excluded")
	}
	if stats.TotalLines == 0 {
		t.Errorf("expected non-zero total lines")
	}
}

func TestSliceWithProgressNilReporter(t *testing.T) {
	from := mustTime("2024-01-10T10:00:00Z")
	to := mustTime("2024-01-10T10:20:00Z")

	src := strings.NewReader(buildProgressLog())
	var dst bytes.Buffer

	_, err := SliceWithProgress(src, &dst, from, to, nil)
	if err != nil {
		t.Fatalf("unexpected error with nil reporter: %v", err)
	}
}

func TestSliceWithProgressEmptyInput(t *testing.T) {
	var dst bytes.Buffer
	var progBuf bytes.Buffer
	reporter := NewProgressReporter(&progBuf, 100)

	_, err := SliceWithProgress(strings.NewReader(""), &dst, time.Now(), time.Now(), reporter)
	if err != nil {
		t.Fatalf("unexpected error on empty input: %v", err)
	}
	if dst.Len() != 0 {
		t.Errorf("expected empty output for empty input")
	}
}
