package slicer

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
	"time"
)

func buildGrepLog() string {
	return strings.Join([]string{
		`2024-01-01T10:00:00Z INFO  service started`,
		`2024-01-01T10:01:00Z ERROR disk full`,
		`2024-01-01T10:02:00Z INFO  request ok`,
		`2024-01-01T10:03:00Z ERROR timeout reached`,
		`2024-01-01T10:04:00Z INFO  shutting down`,
	}, "\n") + "\n"
}

func mustCompile(t *testing.T, pat string) *regexp.Regexp {
	t.Helper()
	return regexp.MustCompile(pat)
}

func TestGrepSliceMatchPattern(t *testing.T) {
	from := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 10, 5, 0, 0, time.UTC)

	var out bytes.Buffer
	err := GrepSlice(strings.NewReader(buildGrepLog()), from, to, GrepOptions{
		Pattern: mustCompile(t, "ERROR"),
	}, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := out.String()
	if !strings.Contains(result, "disk full") {
		t.Errorf("expected 'disk full' in output, got: %s", result)
	}
	if !strings.Contains(result, "timeout reached") {
		t.Errorf("expected 'timeout reached' in output, got: %s", result)
	}
	if strings.Contains(result, "service started") {
		t.Errorf("unexpected 'service started' in output")
	}
}

func TestGrepSliceInvert(t *testing.T) {
	from := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 10, 5, 0, 0, time.UTC)

	var out bytes.Buffer
	err := GrepSlice(strings.NewReader(buildGrepLog()), from, to, GrepOptions{
		Pattern: mustCompile(t, "ERROR"),
		Invert:  true,
	}, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := out.String()
	if strings.Contains(result, "disk full") {
		t.Errorf("unexpected ERROR line in inverted output")
	}
	if !strings.Contains(result, "service started") {
		t.Errorf("expected non-ERROR lines in inverted output")
	}
}

func TestGrepSliceContext(t *testing.T) {
	from := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 10, 5, 0, 0, time.UTC)

	var out bytes.Buffer
	err := GrepSlice(strings.NewReader(buildGrepLog()), from, to, GrepOptions{
		Pattern: mustCompile(t, "disk full"),
		Context: 1,
	}, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := out.String()
	// context=1 should pull in the line before (service started) and after (request ok)
	if !strings.Contains(result, "service started") {
		t.Errorf("expected context line 'service started'")
	}
	if !strings.Contains(result, "request ok") {
		t.Errorf("expected context line 'request ok'")
	}
}

func TestGrepSliceNilPattern(t *testing.T) {
	from := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 10, 5, 0, 0, time.UTC)

	var out bytes.Buffer
	err := GrepSlice(strings.NewReader(buildGrepLog()), from, to, GrepOptions{}, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Len() == 0 {
		t.Error("expected all in-range lines when pattern is nil")
	}
}

func TestGrepSliceNoMatch(t *testing.T) {
	from := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 10, 5, 0, 0, time.UTC)

	var out bytes.Buffer
	err := GrepSlice(strings.NewReader(buildGrepLog()), from, to, GrepOptions{
		Pattern: mustCompile(t, "CRITICAL"),
	}, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Len() != 0 {
		t.Errorf("expected empty output for no-match pattern, got: %s", out.String())
	}
}
