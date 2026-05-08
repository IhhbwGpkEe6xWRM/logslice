package slicer

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func buildUniqueLog(entries []string) *strings.Reader {
	return strings.NewReader(strings.Join(entries, "\n") + "\n")
}

func mustUniqueTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSliceUniqueByMessage(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z INFO disk full",
		"2024-01-01T10:01:00Z INFO disk full",
		"2024-01-01T10:02:00Z INFO connection timeout",
		"2024-01-01T10:03:00Z WARN disk full",
	}

	r := buildUniqueLog(lines)
	var w bytes.Buffer
	n, err := SliceUnique(r, &w, UniqueOptions{
		From:  mustUniqueTime("2024-01-01T09:00:00Z"),
		To:    mustUniqueTime("2024-01-01T11:00:00Z"),
		Field: UniqueByMessage,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 unique messages, got %d", n)
	}
	out := w.String()
	if !strings.Contains(out, "disk full") {
		t.Error("expected 'disk full' in output")
	}
	if !strings.Contains(out, "connection timeout") {
		t.Error("expected 'connection timeout' in output")
	}
}

func TestSliceUniqueByLevel(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z INFO first info",
		"2024-01-01T10:01:00Z INFO second info",
		"2024-01-01T10:02:00Z WARN first warn",
		"2024-01-01T10:03:00Z WARN second warn",
		"2024-01-01T10:04:00Z ERROR something failed",
	}

	r := buildUniqueLog(lines)
	var w bytes.Buffer
	n, err := SliceUnique(r, &w, UniqueOptions{
		From:  mustUniqueTime("2024-01-01T09:00:00Z"),
		To:    mustUniqueTime("2024-01-01T11:00:00Z"),
		Field: UniqueByLevel,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 unique levels, got %d", n)
	}
}

func TestSliceUniqueEmptyInput(t *testing.T) {
	r := strings.NewReader("")
	var w bytes.Buffer
	n, err := SliceUnique(r, &w, UniqueOptions{
		From:  mustUniqueTime("2024-01-01T09:00:00Z"),
		To:    mustUniqueTime("2024-01-01T11:00:00Z"),
		Field: UniqueByMessage,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
}

func TestSliceUniqueOutOfRange(t *testing.T) {
	lines := []string{
		"2024-01-01T08:00:00Z INFO early line",
		"2024-01-01T12:00:00Z INFO late line",
	}

	r := buildUniqueLog(lines)
	var w bytes.Buffer
	n, err := SliceUnique(r, &w, UniqueOptions{
		From:  mustUniqueTime("2024-01-01T09:00:00Z"),
		To:    mustUniqueTime("2024-01-01T11:00:00Z"),
		Field: UniqueByMessage,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
}
