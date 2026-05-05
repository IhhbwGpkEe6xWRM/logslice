package slicer

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func buildMergeLog(lines []string) *strings.Reader {
	return strings.NewReader(strings.Join(lines, "\n") + "\n")
}

func TestMergeSliceSingleReader(t *testing.T) {
	r := buildMergeLog([]string{
		"2024-01-01T10:00:00Z info alpha",
		"2024-01-01T10:01:00Z info beta",
		"2024-01-01T10:02:00Z info gamma",
	})
	from := mustTime("2024-01-01T10:00:00Z")
	to := mustTime("2024-01-01T10:01:30Z")
	var out bytes.Buffer
	n, err := MergeSlice([]interface{ Read([]byte) (int, error) }{r}, from, to, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 lines, got %d", n)
	}
	if !strings.Contains(out.String(), "alpha") {
		t.Error("expected alpha in output")
	}
	if strings.Contains(out.String(), "gamma") {
		t.Error("did not expect gamma in output")
	}
}

func TestMergeSliceMultipleReadersMergesChronologically(t *testing.T) {
	r1 := buildMergeLog([]string{
		"2024-01-01T10:00:00Z info from-r1-first",
		"2024-01-01T10:02:00Z info from-r1-third",
	})
	r2 := buildMergeLog([]string{
		"2024-01-01T10:01:00Z info from-r2-second",
		"2024-01-01T10:03:00Z info from-r2-fourth",
	})
	from := mustTime("2024-01-01T10:00:00Z")
	to := mustTime("2024-01-01T10:03:00Z")
	var out bytes.Buffer
	n, err := MergeSlice([]interface{ Read([]byte) (int, error) }{r1, r2}, from, to, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 4 {
		t.Errorf("expected 4 lines, got %d", n)
	}
	result := out.String()
	posFirst := strings.Index(result, "from-r1-first")
	posSecond := strings.Index(result, "from-r2-second")
	posThird := strings.Index(result, "from-r1-third")
	if posFirst > posSecond || posSecond > posThird {
		t.Errorf("lines not in chronological order: %s", result)
	}
}

func TestMergeSliceEmptyReaders(t *testing.T) {
	r := strings.NewReader("")
	from := mustTime("2024-01-01T10:00:00Z")
	to := mustTime("2024-01-01T11:00:00Z")
	var out bytes.Buffer
	n, err := MergeSlice([]interface{ Read([]byte) (int, error) }{r}, from, to, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
}

func TestMergeSliceNoReadersInRange(t *testing.T) {
	r := buildMergeLog([]string{
		"2024-01-01T08:00:00Z info early",
	})
	from := mustTime("2024-01-01T10:00:00Z")
	to := mustTime("2024-01-01T11:00:00Z")
	var out bytes.Buffer
	n, _ := MergeSlice([]interface{ Read([]byte) (int, error) }{r}, from, to, &out)
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

var _ = time.RFC3339 // suppress unused import
