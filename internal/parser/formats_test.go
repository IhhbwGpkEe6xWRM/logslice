package parser_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func TestTryParseTimestampPrefixRFC3339(t *testing.T) {
	line := "2024-01-15T10:30:00Z some log message here"
	ts, rest, ok := parser.TryParseTimestampPrefix(line)
	if !ok {
		t.Fatal("expected ok=true for RFC3339 prefix")
	}
	if rest != "some log message here" {
		t.Errorf("unexpected rest: %q", rest)
	}
	expected := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if !ts.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, ts)
	}
}

func TestTryParseTimestampPrefixSpaceSeparated(t *testing.T) {
	line := "2024-01-15 10:30:00 some log message"
	ts, rest, ok := parser.TryParseTimestampPrefix(line)
	if !ok {
		t.Fatal("expected ok=true for space-separated prefix")
	}
	if rest != "some log message" {
		t.Errorf("unexpected rest: %q", rest)
	}
	expected := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if !ts.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, ts)
	}
}

func TestTryParseTimestampPrefixNoTimestamp(t *testing.T) {
	line := "no timestamp here at all"
	_, _, ok := parser.TryParseTimestampPrefix(line)
	if ok {
		t.Error("expected ok=false for line with no timestamp")
	}
}

func TestTryParseTimestampPrefixEmptyString(t *testing.T) {
	_, _, ok := parser.TryParseTimestampPrefix("")
	if ok {
		t.Error("expected ok=false for empty string")
	}
}

func TestTryParseTimestampPrefixWithMilliseconds(t *testing.T) {
	line := "2024-03-22T08:15:30.123Z INFO application started"
	ts, rest, ok := parser.TryParseTimestampPrefix(line)
	if !ok {
		t.Fatal("expected ok=true for RFC3339 with milliseconds")
	}
	if rest != "INFO application started" {
		t.Errorf("unexpected rest: %q", rest)
	}
	if ts.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestTryParseTimestampPrefixWithOffset(t *testing.T) {
	line := "2024-06-01T12:00:00+05:30 request received"
	ts, rest, ok := parser.TryParseTimestampPrefix(line)
	if !ok {
		t.Fatal("expected ok=true for RFC3339 with timezone offset")
	}
	if rest != "request received" {
		t.Errorf("unexpected rest: %q", rest)
	}
	if ts.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
