package slicer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/logslice/logslice/internal/slicer"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), "test.log")
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp log: %v", err)
	}
	return tmp
}

func TestSliceFileSuccess(t *testing.T) {
	src := writeTempLog(t, sampleLog)

	var out strings.Builder
	count, err := slicer.SliceFile(src, &out, slicer.Options{
		From: mustTime("2024-01-15T10:02:00Z"),
		To:   mustTime("2024-01-15T10:04:00Z"),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 lines, got %d", count)
	}
}

func TestSliceFileMissing(t *testing.T) {
	_, err := slicer.SliceFile("/nonexistent/file.log", &strings.Builder{}, slicer.Options{
		From: mustTime("2024-01-15T10:00:00Z"),
		To:   mustTime("2024-01-15T11:00:00Z"),
	})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestSliceFileToFile(t *testing.T) {
	src := writeTempLog(t, sampleLog)
	dst := filepath.Join(t.TempDir(), "out.log")

	count, err := slicer.SliceFileToFile(src, dst, slicer.Options{
		From: mustTime("2024-01-15T10:00:00Z"),
		To:   mustTime("2024-01-15T10:01:00Z"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 lines, got %d", count)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !strings.Contains(string(data), "starting server") {
		t.Error("expected output to contain 'starting server'")
	}
}
