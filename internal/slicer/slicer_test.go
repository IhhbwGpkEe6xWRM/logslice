package slicer_test

import (
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/slicer"
)

const sampleLog = `2024-01-15T10:00:00Z INFO starting server
2024-01-15T10:01:00Z INFO listening on :8080
2024-01-15T10:02:00Z WARN high memory usage
2024-01-15T10:03:00Z ERROR connection refused
2024-01-15T10:04:00Z INFO request completed
`

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSliceExtractsRange(t *testing.T) {
	from := mustTime("2024-01-15T10:01:00Z")
	to := mustTime("2024-01-15T10:03:00Z")

	s := slicer.New(slicer.Options{From: from, To: to})
	var out strings.Builder
	count, err := s.Slice(strings.NewReader(sampleLog), &out)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 lines, got %d", count)
	}
	if !strings.Contains(out.String(), "listening on :8080") {
		t.Error("expected line about listening on :8080")
	}
	if strings.Contains(out.String(), "starting server") {
		t.Error("unexpected line: starting server")
	}
}

func TestSliceEmptyInput(t *testing.T) {
	s := slicer.New(slicer.Options{
		From: mustTime("2024-01-15T10:00:00Z"),
		To:   mustTime("2024-01-15T11:00:00Z"),
	})
	var out strings.Builder
	count, err := s.Slice(strings.NewReader(""), &out)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 lines, got %d", count)
	}
}

func TestSliceNoMatchingLines(t *testing.T) {
	s := slicer.New(slicer.Options{
		From: mustTime("2024-01-15T12:00:00Z"),
		To:   mustTime("2024-01-15T13:00:00Z"),
	})
	var out strings.Builder
	count, err := s.Slice(strings.NewReader(sampleLog), &out)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 lines, got %d", count)
	}
}
