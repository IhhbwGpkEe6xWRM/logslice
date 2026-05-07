package slicer

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func writeTempGrepLog(t *testing.T) string {
	t.Helper()
	content := strings.Join([]string{
		`2024-03-01T08:00:00Z INFO  boot complete`,
		`2024-03-01T08:01:00Z ERROR connection refused`,
		`2024-03-01T08:02:00Z INFO  retry attempt`,
		`2024-03-01T08:03:00Z ERROR auth failed`,
	}, "\n") + "\n"

	f, err := os.CreateTemp(t.TempDir(), "grep-*.log")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestGrepFileSuccess(t *testing.T) {
	path := writeTempGrepLog(t)
	from := time.Date(2024, 3, 1, 8, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 8, 5, 0, 0, time.UTC)

	var sb strings.Builder
	err := GrepFile(path, from, to, GrepOptions{
		Pattern: regexp.MustCompile("ERROR"),
	}, &sb)
	if err != nil {
		t.Fatalf("GrepFile: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "connection refused") {
		t.Errorf("expected 'connection refused', got: %s", out)
	}
	if strings.Contains(out, "boot complete") {
		t.Errorf("unexpected INFO line in output")
	}
}

func TestGrepFileMissing(t *testing.T) {
	from := time.Date(2024, 3, 1, 8, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 9, 0, 0, 0, time.UTC)
	var sb strings.Builder
	err := GrepFile("/nonexistent/path.log", from, to, GrepOptions{}, &sb)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestGrepFileToFile(t *testing.T) {
	src := writeTempGrepLog(t)
	dst := filepath.Join(t.TempDir(), "out.log")
	from := time.Date(2024, 3, 1, 8, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 8, 5, 0, 0, time.UTC)

	err := GrepFileToFile(src, dst, from, to, GrepOptions{
		Pattern: regexp.MustCompile("auth failed"),
	})
	if err != nil {
		t.Fatalf("GrepFileToFile: %v", err)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(data), "auth failed") {
		t.Errorf("expected 'auth failed' in output file, got: %s", data)
	}
}
