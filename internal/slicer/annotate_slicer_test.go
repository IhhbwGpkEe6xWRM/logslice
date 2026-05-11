package slicer

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func buildAnnotateLog() string {
	return strings.Join([]string{
		"2024-01-01T10:00:00Z INFO starting server",
		"2024-01-01T10:00:01Z INFO listening on :8080",
		"2024-01-01T10:00:02Z WARN high memory usage",
		"2024-01-01T10:00:03Z ERROR connection refused",
		"2024-01-01T10:00:04Z INFO shutdown complete",
	}, "\n")
}

func mustAnnotateTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSliceAnnotatedLineNumbers(t *testing.T) {
	r := strings.NewReader(buildAnnotateLog())
	from := mustAnnotateTime("2024-01-01T10:00:01Z")
	to := mustAnnotateTime("2024-01-01T10:00:03Z")
	var buf bytes.Buffer
	opts := AnnotateOptions{AddLineNumbers: true}
	n, err := SliceAnnotated(r, from, to, opts, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Fatalf("expected 3 matched lines, got %d", n)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if !strings.HasPrefix(lines[0], "[L2] ") {
		t.Errorf("expected line prefix [L2], got %q", lines[0])
	}
}

func TestSliceAnnotatedRelativeTime(t *testing.T) {
	r := strings.NewReader(buildAnnotateLog())
	from := mustAnnotateTime("2024-01-01T10:00:00Z")
	to := mustAnnotateTime("2024-01-01T10:00:02Z")
	var buf bytes.Buffer
	opts := AnnotateOptions{AddRelativeTime: true}
	n, err := SliceAnnotated(r, from, to, opts, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Fatalf("expected 3 lines, got %d", n)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if !strings.HasPrefix(lines[0], "[+0ms] ") {
		t.Errorf("first line should have +0ms, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "[+1000ms] ") {
		t.Errorf("second line should have +1000ms, got %q", lines[1])
	}
}

func TestSliceAnnotatedEmptyInput(t *testing.T) {
	r := strings.NewReader("")
	from := mustAnnotateTime("2024-01-01T10:00:00Z")
	to := mustAnnotateTime("2024-01-01T10:00:05Z")
	var buf bytes.Buffer
	n, err := SliceAnnotated(r, from, to, AnnotateOptions{}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
}

func TestSliceAnnotatedOffsetAndLineNumbers(t *testing.T) {
	r := strings.NewReader(buildAnnotateLog())
	from := mustAnnotateTime("2024-01-01T10:00:00Z")
	to := mustAnnotateTime("2024-01-01T10:00:00Z")
	var buf bytes.Buffer
	opts := AnnotateOptions{AddLineNumbers: true, AddOffset: true}
	n, err := SliceAnnotated(r, from, to, opts, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 line, got %d", n)
	}
	out := buf.String()
	if !strings.Contains(out, "[L1]") || !strings.Contains(out, "[+") {
		t.Errorf("expected both line number and offset annotations, got %q", out)
	}
}
