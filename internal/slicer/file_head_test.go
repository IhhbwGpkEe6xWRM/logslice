package slicer

import (
	"os"
	"strings"
	"testing"
	"time"
)

func writeTempHeadLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "head-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestHeadFileSuccess(t *testing.T) {
	path := writeTempHeadLog(t, buildHeadLog())
	from := mustTime("2024-01-01T10:00:00Z")
	to := mustTime("2024-01-01T11:00:00Z")

	lines, err := HeadFile(path, from, to, HeadOptions{Lines: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestHeadFileMissing(t *testing.T) {
	from := mustTime("2024-01-01T10:00:00Z")
	to := mustTime("2024-01-01T11:00:00Z")

	_, err := HeadFile("/nonexistent/path.log", from, to, HeadOptions{Lines: 5})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestHeadFileToWriter(t *testing.T) {
	path := writeTempHeadLog(t, buildHeadLog())
	from := mustTime("2024-01-01T10:00:00Z")
	to := mustTime("2024-01-01T11:00:00Z")

	var sb strings.Builder
	err := HeadFileToWriter(path, from, to, HeadOptions{Lines: 3}, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := strings.Split(strings.TrimRight(sb.String(), "\n"), "\n")
	if len(got) != 3 {
		t.Fatalf("expected 3 output lines, got %d", len(got))
	}
}

func TestHeadFileToWriterDuration(t *testing.T) {
	path := writeTempHeadLog(t, buildHeadLog())
	from := mustTime("2024-01-01T10:00:00Z")
	to := mustTime("2024-01-01T11:00:00Z")

	var sb strings.Builder
	err := HeadFileToWriter(path, from, to, HeadOptions{Duration: 15 * time.Second}, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := strings.Split(strings.TrimRight(sb.String(), "\n"), "\n")
	// 0s and 10s fit within 15s; 20s does not
	if len(got) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(got))
	}
}
