package parser_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

type tsCase struct {
	input    string
	wantYear int
	wantErr  bool
}

func TestParseTimestamp(t *testing.T) {
	cases := []tsCase{
		{input: "2024-03-15T08:30:00Z", wantYear: 2024},
		{input: "2024-03-15T08:30:00.123456789Z", wantYear: 2024},
		{input: "2024-03-15 08:30:00", wantYear: 2024},
		{input: "2024-03-15 08:30:00.999", wantYear: 2024},
		{input: "2024/03/15 08:30:00", wantYear: 2024},
		{input: "15/Mar/2024:08:30:00 +0000", wantYear: 2024},
		{input: "not-a-timestamp", wantErr: true},
		{input: "", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, layout, err := parser.ParseTimestamp(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for input %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tc.input, err)
			}
			if got.Year() != tc.wantYear {
				t.Errorf("year: got %d, want %d", got.Year(), tc.wantYear)
			}
			if layout == "" {
				t.Error("expected non-empty layout string")
			}
		})
	}
}

func TestMustParseTimestampPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid timestamp")
		}
	}()
	parser.MustParseTimestamp("garbage")
}

func TestMustParseTimestampValid(t *testing.T) {
	raw := "2024-06-01T12:00:00Z"
	got := parser.MustParseTimestamp(raw)
	if got.IsZero() {
		t.Error("expected non-zero time")
	}
	if got.Year() != 2024 || got.Month() != time.June {
		t.Errorf("unexpected time value: %v", got)
	}
}
