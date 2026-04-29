package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunMissingInput(t *testing.T) {
	err := Run([]string{"--from", "2024-01-01T00:00:00Z", "--to", "2024-01-02T00:00:00Z"})
	if err == nil || !strings.Contains(err.Error(), "--input") {
		t.Errorf("expected --input error, got %v", err)
	}
}

func TestRunMissingFrom(t *testing.T) {
	err := Run([]string{"--input", "x.log", "--to", "2024-01-02T00:00:00Z"})
	if err == nil || !strings.Contains(err.Error(), "--from") {
		t.Errorf("expected --from error, got %v", err)
	}
}

func TestRunFromAfterTo(t *testing.T) {
	err := Run([]string{
		"--input", "x.log",
		"--from", "2024-01-02T00:00:00Z",
		"--to", "2024-01-01T00:00:00Z",
	})
	if err == nil || !strings.Contains(err.Error(), "before") {
		t.Errorf("expected ordering error, got %v", err)
	}
}

func TestRunToStdout(t *testing.T) {
	lines := "2024-03-01T10:00:00Z INFO starting\n" +
		"2024-03-01T10:05:00Z INFO running\n" +
		"2024-03-01T11:00:00Z INFO done\n"

	tmpFile := filepath.Join(t.TempDir(), "test.log")
	if err := os.WriteFile(tmpFile, []byte(lines), 0o644); err != nil {
		t.Fatal(err)
	}

	// Redirect stdout via sliceToWriter directly
	var buf bytes.Buffer
	err := sliceToWriter(tmpFile,
		mustParseTime(t, "2024-03-01T09:59:00Z"),
		mustParseTime(t, "2024-03-01T10:30:00Z"),
		&buf,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "starting") {
		t.Errorf("expected 'starting' in output, got: %s", out)
	}
	if !strings.Contains(out, "running") {
		t.Errorf("expected 'running' in output, got: %s", out)
	}
	if strings.Contains(out, "done") {
		t.Errorf("did not expect 'done' in output, got: %s", out)
	}
}

func mustParseTime(t *testing.T, s string) (result interface{}) {
	t.Helper()
	import_time := func() {}
	_ = import_time
	// Use the parser package indirectly via parseArgs
	cfg, err := parseArgs([]string{
		"--input", "dummy",
		"--from", s,
		"--to", "2099-01-01T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("mustParseTime(%q): %v", s, err)
	}
	return cfg.From
}
