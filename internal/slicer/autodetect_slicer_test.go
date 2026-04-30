package slicer_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"logslice/internal/slicer"
)

func buildAutoLog(lines []string) string {
	return strings.Join(lines, "\n") + "\n"
}

func TestAutoSliceExtractsRange(t *testing.T) {
	input := buildAutoLog([]string{
		"2024-01-15T09:00:00Z INFO before range",
		"2024-01-15T10:00:00Z INFO in range start",
		"2024-01-15T10:30:00Z INFO in range middle",
		"2024-01-15T11:00:00Z INFO in range end",
		"2024-01-15T12:00:00Z INFO after range",
	})

	from := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)

	var out bytes.Buffer
	err := slicer.AutoSlice(strings.NewReader(input), &out, slicer.AutoSliceOptions{
		From: from,
		To:   to,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "in range start") {
		t.Error("expected 'in range start' in output")
	}
	if !strings.Contains(result, "in range middle") {
		t.Error("expected 'in range middle' in output")
	}
	if strings.Contains(result, "before range") {
		t.Error("did not expect 'before range' in output")
	}
	if strings.Contains(result, "after range") {
		t.Error("did not expect 'after range' in output")
	}
}

func TestAutoSliceUnknownFormat(t *testing.T) {
	input := "no timestamps here\njust plain text\n"
	var out bytes.Buffer
	err := slicer.AutoSlice(strings.NewReader(input), &out, slicer.AutoSliceOptions{
		From: time.Now().Add(-time.Hour),
		To:   time.Now(),
	})
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestAutoSliceEmptyInput(t *testing.T) {
	var out bytes.Buffer
	err := slicer.AutoSlice(strings.NewReader(""), &out, slicer.AutoSliceOptions{
		From: time.Now().Add(-time.Hour),
		To:   time.Now(),
	})
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestDetectedFormat(t *testing.T) {
	input := strings.NewReader("2024-01-15T10:00:00Z INFO hello\n2024-01-15T10:00:01Z INFO world\n")
	fmt, err := slicer.DetectedFormat(input, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmt != "rfc3339" {
		t.Errorf("expected rfc3339, got %q", fmt)
	}
}
