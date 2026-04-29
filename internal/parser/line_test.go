package parser

import (
	"testing"
	"time"
)

func TestParseLineRFC3339(t *testing.T) {
	raw := "2024-03-15T12:34:56Z [INFO] server started"
	line, err := ParseLine(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if line.Raw != raw {
		t.Errorf("Raw mismatch: got %q want %q", line.Raw, raw)
	}
	if line.Timestamp.Year() != 2024 || line.Timestamp.Month() != 3 || line.Timestamp.Day() != 15 {
		t.Errorf("unexpected timestamp: %v", line.Timestamp)
	}
}

func TestParseLineSpaceSeparated(t *testing.T) {
	raw := "2024-03-15 08:00:01 [WARN] disk usage high"
	line, err := ParseLine(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if line.Timestamp.Hour() != 8 {
		t.Errorf("expected hour 8, got %d", line.Timestamp.Hour())
	}
}

func TestParseLineNoTimestamp(t *testing.T) {
	_, err := ParseLine("this line has no timestamp at all")
	if err != ErrNoTimestamp {
		t.Errorf("expected ErrNoTimestamp, got %v", err)
	}
}

func TestParseLineEmptyString(t *testing.T) {
	_, err := ParseLine("")
	if err != ErrNoTimestamp {
		t.Errorf("expected ErrNoTimestamp for empty string, got %v", err)
	}
}

func TestLogLineInRange(t *testing.T) {
	base := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)
	line := LogLine{Timestamp: base}

	start := base.Add(-time.Hour)
	end := base.Add(time.Hour)

	if !line.InRange(start, end) {
		t.Error("expected line to be in range")
	}
	if line.InRange(base.Add(time.Minute), end) {
		t.Error("expected line to be out of range (before start)")
	}
	if line.InRange(start, base.Add(-time.Minute)) {
		t.Error("expected line to be out of range (after end)")
	}
}

func TestLogLineInRangeOpenBounds(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	line := LogLine{Timestamp: base}

	// Both bounds open.
	if !line.InRange(time.Time{}, time.Time{}) {
		t.Error("expected line to be in range with open bounds")
	}
	// Only start bound.
	if !line.InRange(base.Add(-time.Second), time.Time{}) {
		t.Error("expected line to be in range with open end")
	}
	// Only end bound.
	if !line.InRange(time.Time{}, base.Add(time.Second)) {
		t.Error("expected line to be in range with open start")
	}
}
