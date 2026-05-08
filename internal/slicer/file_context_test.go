package slicer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeTempContextLog(t *testing.T) string {
	t.Helper()
	content := strings.Join([]string{
		"2024-03-01T08:00:00Z INFO  startup",
		"2024-03-01T08:01:00Z WARN  low memory",
		"2024-03-01T08:02:00Z ERROR disk full",
		"2024-03-01T08:03:00Z INFO  retrying",
		"2024-03-01T08:04:00Z INFO  recovered",
	}, "\n")
	f, err := os.CreateTemp(t.TempDir(), "ctx*.log")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func mustContextFileTime(s string) time.Time {
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return v
}

func TestContextFileSuccess(t *testing.T) {
	path := writeTempContextLog(t)
	from := mustContextFileTime("2024-03-01T08:02:00Z")
	to := mustContextFileTime("2024-03-01T08:02:00Z")

	var sb strings.Builder
	err := ContextFile(path, from, to, ContextOptions{Before: 1, After: 1}, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "low memory") {
		t.Errorf("expected before-context line, got:\n%s", out)
	}
	if !strings.Contains(out, "disk full") {
		t.Errorf("expected matched line, got:\n%s", out)
	}
	if !strings.Contains(out, "retrying") {
		t.Errorf("expected after-context line, got:\n%s", out)
	}
}

func TestContextFileMissing(t *testing.T) {
	var sb strings.Builder
	err := ContextFile("/nonexistent/file.log",
		mustContextFileTime("2024-03-01T08:00:00Z"),
		mustContextFileTime("2024-03-01T09:00:00Z"),
		ContextOptions{}, &sb)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestContextFileToFile(t *testing.T) {
	src := writeTempContextLog(t)
	dst := filepath.Join(t.TempDir(), "out.log")

	from := mustContextFileTime("2024-03-01T08:01:00Z")
	to := mustContextFileTime("2024-03-01T08:03:00Z")

	if err := ContextFileToFile(src, dst, from, to, ContextOptions{Before: 1}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("could not read output: %v", err)
	}
	if !strings.Contains(string(data), "startup") {
		t.Errorf("expected before-context 'startup' in output:\n%s", data)
	}
	if !strings.Contains(string(data), "retrying") {
		t.Errorf("expected matched line 'retrying' in output:\n%s", data)
	}
}
