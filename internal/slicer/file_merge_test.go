package slicer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempMergeLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "merge-*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	defer f.Close()
	f.WriteString(strings.Join(lines, "\n") + "\n")
	return f.Name()
}

func TestMergeFilesSuccess(t *testing.T) {
	p1 := writeTempMergeLog(t, []string{
		"2024-03-01T09:00:00Z info service-a started",
		"2024-03-01T09:02:00Z info service-a ready",
	})
	p2 := writeTempMergeLog(t, []string{
		"2024-03-01T09:01:00Z info service-b connected",
		"2024-03-01T09:03:00Z info service-b synced",
	})
	from := mustTime("2024-03-01T09:00:00Z")
	to := mustTime("2024-03-01T09:03:00Z")
	var sb strings.Builder
	n, err := MergeFiles([]string{p1, p2}, from, to, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 4 {
		t.Errorf("expected 4 lines, got %d", n)
	}
	if !strings.Contains(sb.String(), "service-a") || !strings.Contains(sb.String(), "service-b") {
		t.Error("expected both services in output")
	}
}

func TestMergeFilesMissing(t *testing.T) {
	_, err := MergeFiles([]string{"/nonexistent/file.log"}, mustTime("2024-01-01T00:00:00Z"), mustTime("2024-01-02T00:00:00Z"), &strings.Builder{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestMergeFilesNoInputs(t *testing.T) {
	_, err := MergeFiles([]string{}, mustTime("2024-01-01T00:00:00Z"), mustTime("2024-01-02T00:00:00Z"), &strings.Builder{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestMergeFilesToFile(t *testing.T) {
	p1 := writeTempMergeLog(t, []string{
		"2024-03-01T10:00:00Z info alpha",
	})
	out := filepath.Join(t.TempDir(), "merged.log")
	from := mustTime("2024-03-01T10:00:00Z")
	to := mustTime("2024-03-01T10:01:00Z")
	n, err := MergeFilesToFile([]string{p1}, from, to, out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 line, got %d", n)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "alpha") {
		t.Error("expected alpha in output file")
	}
}

func TestMergeFilesOutputOrdering(t *testing.T) {
	p1 := writeTempMergeLog(t, []string{
		"2024-03-01T09:00:00Z info first",
		"2024-03-01T09:02:00Z info third",
	})
	p2 := writeTempMergeLog(t, []string{
		"2024-03-01T09:01:00Z info second",
		"2024-03-01T09:03:00Z info fourth",
	})
	from := mustTime("2024-03-01T09:00:00Z")
	to := mustTime("2024-03-01T09:03:00Z")
	var sb strings.Builder
	_, err := MergeFiles([]string{p1, p2}, from, to, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := sb.String()
	posFirst := strings.Index(output, "first")
	posSecond := strings.Index(output, "second")
	posThird := strings.Index(output, "third")
	if posFirst > posSecond || posSecond > posThird {
		t.Error("expected output lines to be in chronological order")
	}
}
