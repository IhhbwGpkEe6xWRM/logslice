package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSplitInput(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "split-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	f.WriteString(strings.Join(lines, "\n") + "\n")
	return f.Name()
}

func TestRunSplitMissingArgs(t *testing.T) {
	err := RunSplit([]string{"only-one"}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunSplitMissingInput(t *testing.T) {
	err := RunSplit([]string{"/no/such/file.log", t.TempDir(), "1h"}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for missing input")
	}
}

func TestRunSplitInvalidWindow(t *testing.T) {
	input := writeSplitInput(t, []string{"2024-01-01T00:00:00Z level=info msg=x"})
	err := RunSplit([]string{input, t.TempDir(), "notaduration"}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for invalid window")
	}
}

func TestRunSplitSuccess(t *testing.T) {
	lines := []string{
		"2024-06-01T10:00:00Z level=info msg=a",
		"2024-06-01T10:30:00Z level=info msg=b",
		"2024-06-01T11:05:00Z level=info msg=c",
	}
	input := writeSplitInput(t, lines)
	outDir := t.TempDir()
	var stderr bytes.Buffer

	err := RunSplit([]string{input, outDir, "1h"}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stderr.String(), "split complete") {
		t.Errorf("expected completion message, got: %s", stderr.String())
	}

	entries, _ := os.ReadDir(outDir)
	if len(entries) < 1 {
		t.Error("expected at least one output file")
	}
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".log" {
			t.Errorf("unexpected file extension: %s", e.Name())
		}
	}
}

func TestParseSplitWindowDuration(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"1h", false},
		{"30m", false},
		{"60", false},
		{"-1h", true},
		{"0", true},
		{"bad", true},
	}
	for _, tc := range cases {
		_, err := parseSplitWindow(tc.input)
		if (err != nil) != tc.wantErr {
			t.Errorf("parseSplitWindow(%q) error=%v, wantErr=%v", tc.input, err, tc.wantErr)
		}
	}
}
