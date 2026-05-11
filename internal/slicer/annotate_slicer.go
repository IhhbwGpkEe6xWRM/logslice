package slicer

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// AnnotateOptions controls how lines are annotated.
type AnnotateOptions struct {
	AddLineNumbers bool
	AddRelativeTime bool
	AddOffset bool
	Baseline time.Time
}

// AnnotatedLine wraps a parsed log line with annotation metadata.
type AnnotatedLine struct {
	parser.LogLine
	LineNumber  int
	RelativeMs  int64
	ByteOffset  int64
	Annotated   string
}

// SliceAnnotated reads lines from r within [from, to] and writes annotated
// output to w according to opts.
func SliceAnnotated(r io.Reader, from, to time.Time, opts AnnotateOptions, w io.Writer) (int, error) {
	scanner := bufio.NewScanner(r)
	lineNum := 0
	matched := 0
	var byteOffset int64
	baseline := opts.Baseline

	for scanner.Scan() {
		raw := scanner.Text()
		lineNum++
		byteOffset += int64(len(raw)) + 1

		ll := parser.ParseLine(raw)
		if ll.Timestamp.IsZero() {
			continue
		}
		if ll.Timestamp.Before(from) {
			continue
		}
		if ll.Timestamp.After(to) {
			break
		}

		if opts.AddRelativeTime && baseline.IsZero() {
			baseline = ll.Timestamp
		}

		annotated := buildAnnotation(ll, lineNum, byteOffset, baseline, opts)
		if _, err := fmt.Fprintln(w, annotated); err != nil {
			return matched, err
		}
		matched++
	}
	return matched, scanner.Err()
}

func buildAnnotation(ll parser.LogLine, lineNum int, offset int64, baseline time.Time, opts AnnotateOptions) string {
	prefix := ""
	if opts.AddLineNumbers {
		prefix += fmt.Sprintf("[L%d] ", lineNum)
	}
	if opts.AddOffset {
		prefix += fmt.Sprintf("[+%db] ", offset)
	}
	if opts.AddRelativeTime && !baseline.IsZero() {
		relMs := ll.Timestamp.Sub(baseline).Milliseconds()
		prefix += fmt.Sprintf("[+%dms] ", relMs)
	}
	return prefix + ll.Raw
}
