package parser

import (
	"testing"
	"time"
)

func TestNewLineFilterNoLevelNoPattern(t *testing.T) {
	f, err := NewLineFilter("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match("anything goes here", time.Time{}) {
		t.Error("expected match for empty filter")
	}
}

func TestNewLineFilterInvalidPattern(t *testing.T) {
	_, err := NewLineFilter("", "[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regexp")
	}
}

func TestLineFilterByLevel(t *testing.T) {
	f, err := NewLineFilter("error", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := []struct {
		line string
		want bool
	}{
		{"2024-01-01T00:00:00Z ERROR something failed", true},
		{"2024-01-01T00:00:00Z error lowercase", true},
		{"2024-01-01T00:00:00Z WARN not an error", false},
		{"2024-01-01T00:00:00Z INFO startup", false},
	}
	for _, tc := range cases {
		got := f.Match(tc.line, time.Time{})
		if got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}

func TestLineFilterByPattern(t *testing.T) {
	f, err := NewLineFilter("", `user_id=\d+`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !f.Match("request user_id=42 completed", time.Time{}) {
		t.Error("expected match for pattern")
	}
	if f.Match("request completed without user", time.Time{}) {
		t.Error("expected no match for pattern")
	}
}

func TestLineFilterCombined(t *testing.T) {
	f, err := NewLineFilter("WARN", `timeout`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !f.Match("WARN connection timeout exceeded", time.Time{}) {
		t.Error("expected match for combined filter")
	}
	if f.Match("WARN something else", time.Time{}) {
		t.Error("expected no match: pattern not satisfied")
	}
	if f.Match("ERROR timeout occurred", time.Time{}) {
		t.Error("expected no match: level not satisfied")
	}
}

func TestNilFilterAlwaysMatches(t *testing.T) {
	var f *LineFilter
	if !f.Match("any line", time.Time{}) {
		t.Error("nil filter should always match")
	}
}
