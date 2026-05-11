package slicer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeTempAnnotateLog(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "annotate.log")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempAnnotateLog: %v", err)
	}
	return path
}

func mustAnnotateFileTime(s string) time.Time {
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return v
}

func TestAnnotateFileSuccess(t *testing.T) {
	path := writeTempAnnotateLog(t, buildAnnotateLog())
	from := mustAnnotateFileTime("2024-01-01T10:00:01Z")
	to := mustAnnotateFileTime("2024-01-01T10:00:02Z")
	opts := AnnotateOptions{AddLineNumbers: true, AddRelativeTime: true}
	var buf strings.Builder
	n, err := AnnotateFile(path, from, to, opts, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 lines, got %d", n)
	}
	if !strings.Contains(buf.String(), "[L2]") {
		t.Errorf("expected line number annotation in output")
	}
}

func TestAnnotateFileMissing(t *testing.T) {
	from := mustAnnotateFileTime("2024-01-01T10:00:00Z")
	to := mustAnnotateFileTime("2024-01-01T10:00:05Z")
	var buf strings.Builder
	_, err := AnnotateFile("/nonexistent/annotate.log", from, to, AnnotateOptions{}, &buf)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestAnnotateFileToFile(t *testing.T) {
	src := writeTempAnnotateLog(t, buildAnnotateLog())
	dst := filepath.Join(t.TempDir(), "out.log")
	from := mustAnnotateFileTime("2024-01-01T10:00:00Z")
	to := mustAnnotateFileTime("2024-01-01T10:00:04Z")
	opts := AnnotateOptions{AddLineNumbers: true, AddOffset: true}
	n, err := AnnotateFileToFile(src, dst, from, to, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("expected 5 lines, got %d", n)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if !strings.Contains(string(data), "[L1]") {
		t.Errorf("expected annotation in output file")
	}
}
