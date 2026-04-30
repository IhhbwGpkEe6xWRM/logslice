package slicer

import (
	"fmt"
	"io"
	"sync/atomic"
)

// ProgressReporter reports slicing progress to a writer.
type ProgressReporter struct {
	total     int64
	processed atomic.Int64
	matched   atomic.Int64
	out       io.Writer
	enabled   bool
}

// NewProgressReporter creates a ProgressReporter. If out is nil or total <= 0,
// progress reporting is disabled.
func NewProgressReporter(out io.Writer, totalBytes int64) *ProgressReporter {
	return &ProgressReporter{
		total:   totalBytes,
		out:     out,
		enabled: out != nil && totalBytes > 0,
	}
}

// RecordLine records a processed line of the given byte length and whether it matched.
func (p *ProgressReporter) RecordLine(byteLen int, matched bool) {
	if !p.enabled {
		return
	}
	p.processed.Add(int64(byteLen))
	if matched {
		p.matched.Add(1)
	}
}

// Report writes the current progress to the output writer.
func (p *ProgressReporter) Report() {
	if !p.enabled {
		return
	}
	processed := p.processed.Load()
	pct := float64(processed) / float64(p.total) * 100
	matched := p.matched.Load()
	fmt.Fprintf(p.out, "\rprogress: %.1f%% (%d bytes processed, %d lines matched)",
		pct, processed, matched)
}

// Finish writes a final newline to complete the progress line.
func (p *ProgressReporter) Finish() {
	if !p.enabled {
		return
	}
	p.Report()
	fmt.Fprintln(p.out)
}
