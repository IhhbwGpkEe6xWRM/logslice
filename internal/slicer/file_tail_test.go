package slicer

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

func writeTempTailLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "logslice-tail-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestTailFileSuccess(t *testing.T) {
	content := strings.Join([]string{
		"2024-03-01T09:00:00Z info alpha",
		"2024-03-01T09:01:00Z info beta",
		"2024-03-01T09:02:00Z info gamma",
	}, "\n")
	path := writeTempTailLog(t, content)

	res, err := TailFile(path, TailOptions{Lines: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(res.Lines))
	}
	if !strings.Contains(res.Lines[0], "beta") {
		t.Errorf("expected beta in first result line, got %q", res.Lines[0])
	}
}

func TestTailFileMissing(t *testing.T) {
	_, err := TailFile("/nonexistent/path/logfile.log", TailOptions{Lines: 5})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestTailFileToWriter(t *testing.T) {
	content := strings.Join([]string{
		"2024-03-01T09:00:00Z info one",
		"2024-03-01T09:01:00Z info two",
		"2024-03-01T09:02:00Z info three",
	}, "\n")
	path := writeTempTailLog(t, content)

	var buf bytes.Buffer
	stats, err := TailFileToWriter(path, TailOptions{Duration: 90 * time.Second}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.MatchedLines != 2 {
		t.Errorf("expected 2 matched lines, got %d", stats.MatchedLines)
	}
	output := buf.String()
	if !strings.Contains(output, "two") || !strings.Contains(output, "three") {
		t.Errorf("unexpected output: %q", output)
	}
}
