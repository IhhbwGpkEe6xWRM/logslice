package parser

import (
	"strings"
	"testing"
)

func TestDetectFormatRFC3339(t *testing.T) {
	input := strings.NewReader(
		"2024-01-15T08:00:00Z INFO starting server\n" +
			"2024-01-15T08:00:01Z DEBUG connection accepted\n" +
			"2024-01-15T08:00:02Z INFO request handled\n",
	)
	got := DetectFormat(input, 20)
	if got != FormatRFC3339 {
		t.Errorf("expected RFC3339, got %s", got)
	}
}

func TestDetectFormatSpaceSeparated(t *testing.T) {
	input := strings.NewReader(
		"2024-01-15 08:00:00 INFO starting server\n" +
			"2024-01-15 08:00:01 DEBUG connection accepted\n" +
			"2024-01-15 08:00:02 INFO request handled\n",
	)
	got := DetectFormat(input, 20)
	if got != FormatSpaceSeparated {
		t.Errorf("expected space-separated, got %s", got)
	}
}

func TestDetectFormatMilliseconds(t *testing.T) {
	input := strings.NewReader(
		"1705305600000 INFO starting server\n" +
			"1705305601000 DEBUG connection accepted\n" +
			"1705305602000 INFO request handled\n",
	)
	got := DetectFormat(input, 20)
	if got != FormatMilliseconds {
		t.Errorf("expected milliseconds, got %s", got)
	}
}

func TestDetectFormatUnknown(t *testing.T) {
	input := strings.NewReader(
		"no timestamp here\n" +
			"just plain text\n",
	)
	got := DetectFormat(input, 20)
	if got != FormatUnknown {
		t.Errorf("expected unknown, got %s", got)
	}
}

func TestDetectFormatEmptyInput(t *testing.T) {
	input := strings.NewReader("")
	got := DetectFormat(input, 20)
	if got != FormatUnknown {
		t.Errorf("expected unknown for empty input, got %s", got)
	}
}

func TestDetectFormatDefaultProbeLines(t *testing.T) {
	// maxProbeLines <= 0 should default to 20
	input := strings.NewReader(
		"2024-01-15T08:00:00Z INFO line one\n" +
			"2024-01-15T08:00:01Z INFO line two\n",
	)
	got := DetectFormat(input, 0)
	if got != FormatRFC3339 {
		t.Errorf("expected RFC3339 with default probe lines, got %s", got)
	}
}

func TestFormatString(t *testing.T) {
	cases := []struct {
		f    Format
		want string
	}{
		{FormatRFC3339, "RFC3339"},
		{FormatSpaceSeparated, "space-separated"},
		{FormatMilliseconds, "milliseconds"},
		{FormatUnknown, "unknown"},
	}
	for _, tc := range cases {
		if got := tc.f.String(); got != tc.want {
			t.Errorf("Format(%d).String() = %q, want %q", tc.f, got, tc.want)
		}
	}
}
