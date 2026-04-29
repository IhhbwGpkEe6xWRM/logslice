package cli

import (
	"testing"
	"time"
)

// parseTimeForTest is a test helper that parses RFC3339 or panics.
func parseTimeForTest(t *testing.T, s string) time.Time {
	t.Helper()
	cfg, err := parseArgs([]string{
		"--input", "dummy.log",
		"--from", s,
		"--to", "2099-12-31T23:59:59Z",
	})
	if err != nil {
		t.Fatalf("parseTimeForTest(%q) failed: %v", s, err)
	}
	return cfg.From
}

func TestParseArgsValid(t *testing.T) {
	cfg, err := parseArgs([]string{
		"--input", "app.log",
		"--from", "2024-06-01T08:00:00Z",
		"--to", "2024-06-01T09:00:00Z",
		"--output", "out.log",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Input != "app.log" {
		t.Errorf("expected input=app.log, got %s", cfg.Input)
	}
	if cfg.Output != "out.log" {
		t.Errorf("expected output=out.log, got %s", cfg.Output)
	}
	expectedFrom := time.Date(2024, 6, 1, 8, 0, 0, 0, time.UTC)
	if !cfg.From.Equal(expectedFrom) {
		t.Errorf("expected From=%v, got %v", expectedFrom, cfg.From)
	}
}

func TestParseArgsInvalidFrom(t *testing.T) {
	_, err := parseArgs([]string{
		"--input", "app.log",
		"--from", "not-a-date",
		"--to", "2024-06-01T09:00:00Z",
	})
	if err == nil {
		t.Error("expected error for invalid --from timestamp")
	}
}

func TestParseArgsOutputOptional(t *testing.T) {
	cfg, err := parseArgs([]string{
		"--input", "app.log",
		"--from", "2024-06-01T08:00:00Z",
		"--to", "2024-06-01T09:00:00Z",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Output != "" {
		t.Errorf("expected empty output, got %s", cfg.Output)
	}
}

func TestParseTimeForTest(t *testing.T) {
	ts := parseTimeForTest(t, "2024-01-15T12:30:00Z")
	expected := time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC)
	if !ts.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, ts)
	}
}
