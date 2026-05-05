package slicer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeTempSplitLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "splitlog-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(strings.Join(lines, "\n") + "\n")
	if err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func TestSplitFileByHour(t *testing.T) {
	lines := []string{
		"2024-01-15T10:00:00Z level=info msg=a",
		"2024-01-15T10:30:00Z level=info msg=b",
		"2024-01-15T11:00:00Z level=info msg=c",
		"2024-01-15T11:45:00Z level=info msg=d",
		"2024-01-15T12:00:00Z level=info msg=e",
	}
	src := writeTempSplitLog(t, lines)
	outDir := t.TempDir()

	n, err := SplitFile(src, outDir, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 output files, got %d", n)
	}

	entries, _ := os.ReadDir(outDir)
	if len(entries) != 3 {
		t.Errorf("expected 3 files in outDir, got %d", len(entries))
	}
}

func TestSplitFileMissing(t *testing.T) {
	_, err := SplitFile("/nonexistent/file.log", t.TempDir(), time.Hour)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSplitFileOutputDir(t *testing.T) {
	lines := []string{
		"2024-03-01T08:00:00Z level=info msg=start",
		"2024-03-01T08:59:00Z level=info msg=mid",
		"2024-03-01T09:01:00Z level=info msg=end",
	}
	src := writeTempSplitLog(t, lines)
	outDir := t.TempDir()

	_, err := SplitFile(src, outDir, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := os.ReadDir(outDir)
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".log" {
			t.Errorf("expected .log extension, got %s", e.Name())
		}
	}
}

func TestSplitFileEmptyInput(t *testing.T) {
	src := writeTempSplitLog(t, []string{})
	outDir := t.TempDir()

	n, err := SplitFile(src, outDir, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 output files, got %d", n)
	}
}
