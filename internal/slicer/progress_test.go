package slicer

import (
	"bytes"
	"strings"
	"testing"
)

func TestProgressReporterDisabledWhenNoWriter(t *testing.T) {
	p := NewProgressReporter(nil, 1000)
	p.RecordLine(100, true)
	p.Report() // should not panic
	p.Finish()
}

func TestProgressReporterDisabledWhenZeroTotal(t *testing.T) {
	var buf bytes.Buffer
	p := NewProgressReporter(&buf, 0)
	p.RecordLine(50, true)
	p.Finish()
	if buf.Len() != 0 {
		t.Errorf("expected no output when total=0, got %q", buf.String())
	}
}

func TestProgressReporterOutput(t *testing.T) {
	var buf bytes.Buffer
	p := NewProgressReporter(&buf, 200)
	p.RecordLine(100, true)
	p.RecordLine(100, false)
	p.Finish()

	out := buf.String()
	if !strings.Contains(out, "100.0%") {
		t.Errorf("expected 100.0%% in output, got %q", out)
	}
	if !strings.Contains(out, "1 lines matched") {
		t.Errorf("expected '1 lines matched' in output, got %q", out)
	}
}

func TestProgressReporterPartial(t *testing.T) {
	var buf bytes.Buffer
	p := NewProgressReporter(&buf, 1000)
	p.RecordLine(250, true)
	p.Report()

	out := buf.String()
	if !strings.Contains(out, "25.0%") {
		t.Errorf("expected 25.0%% in output, got %q", out)
	}
}

func TestProgressMatchedCount(t *testing.T) {
	var buf bytes.Buffer
	p := NewProgressReporter(&buf, 500)
	for i := 0; i < 5; i++ {
		p.RecordLine(50, i%2 == 0) // 3 matched (0,2,4)
	}
	p.Finish()

	out := buf.String()
	if !strings.Contains(out, "3 lines matched") {
		t.Errorf("expected '3 lines matched' in output, got %q", out)
	}
}
