package slicer

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeTempRotateLog(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempRotateLog: %v", err)
	}
	return path
}

func TestRotateSliceExtractsRange(t *testing.T) {
	dir := t.TempDir()
	writeTempRotateLog(t, dir, "app.log.2", strings.Join([]string{
		"2024-01-01T08:00:00Z level=info msg=boot",
		"2024-01-01T08:30:00Z level=info msg=ready",
	}, "\n")+"\n")
	writeTempRotateLog(t, dir, "app.log.1", strings.Join([]string{
		"2024-01-01T09:00:00Z level=info msg=request",
		"2024-01-01T09:45:00Z level=warn msg=slow",
	}, "\n")+"\n")
	writeTempRotateLog(t, dir, "app.log", strings.Join([]string{
		"2024-01-01T10:00:00Z level=info msg=done",
	}, "\n")+"\n")

	from := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	var buf bytes.Buffer
	err := RotateSlice(RotateSliceOptions{
		Dir:     dir,
		Pattern: "app.log*",
		From:    from,
		To:      to,
	}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "msg=request") {
		t.Errorf("expected msg=request in output, got:\n%s", got)
	}
	if !strings.Contains(got, "msg=slow") {
		t.Errorf("expected msg=slow in output, got:\n%s", got)
	}
	if strings.Contains(got, "msg=boot") {
		t.Errorf("did not expect msg=boot in output, got:\n%s", got)
	}
	if strings.Contains(got, "msg=done") {
		t.Errorf("did not expect msg=done in output, got:\n%s", got)
	}
}

func TestRotateSliceMissingDir(t *testing.T) {
	var buf bytes.Buffer
	err := RotateSlice(RotateSliceOptions{
		Dir:     "",
		Pattern: "app.log*",
		From:    time.Now(),
	}, &buf)
	if err == nil {
		t.Fatal("expected error for empty dir")
	}
}

func TestRotateSliceNoMatch(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	err := RotateSlice(RotateSliceOptions{
		Dir:     dir,
		Pattern: "nofile*.log",
		From:    time.Now(),
	}, &buf)
	if err == nil {
		t.Fatal("expected error when no files match")
	}
}

func TestRotateSliceFromAfterTo(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	from := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	err := RotateSlice(RotateSliceOptions{
		Dir:     dir,
		Pattern: "app.log*",
		From:    from,
		To:      to,
	}, &buf)
	if err == nil {
		t.Fatal("expected error when from is after to")
	}
}

func TestListRotatedFiles(t *testing.T) {
	dir := t.TempDir()
	writeTempRotateLog(t, dir, "app.log.1", "2024-01-01T08:00:00Z level=info msg=a\n")
	writeTempRotateLog(t, dir, "app.log", "2024-01-01T09:00:00Z level=info msg=b\n")

	files, _, err := ListRotatedFiles(dir, "app.log*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
}
