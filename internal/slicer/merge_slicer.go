package slicer

import (
	"io"
	"sort"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// MergeEntry holds a parsed log line with its source index for merge tracking.
type MergeEntry struct {
	Line      parser.LogLine
	SourceIdx int
}

// MergeSlice reads from multiple readers, merges lines in chronological order,
// and writes lines whose timestamps fall within [from, to] to w.
func MergeSlice(readers []io.Reader, from, to time.Time, w io.Writer) (int, error) {
	var allEntries []MergeEntry

	for idx, r := range readers {
		lines, err := readAllLines(r)
		if err != nil {
			return 0, err
		}
		for _, line := range lines {
			allEntries = append(allEntries, MergeEntry{Line: line, SourceIdx: idx})
		}
	}

	sort.SliceStable(allEntries, func(i, j int) bool {
		if allEntries[i].Line.Timestamp == nil || allEntries[j].Line.Timestamp == nil {
			return false
		}
		return allEntries[i].Line.Timestamp.Before(*allEntries[j].Line.Timestamp)
	})

	written := 0
	for _, entry := range allEntries {
		if entry.Line.InRange(from, to) {
			if _, err := io.WriteString(w, entry.Line.Raw+"\n"); err != nil {
				return written, err
			}
			written++
		}
	}
	return written, nil
}

func readAllLines(r io.Reader) ([]parser.LogLine, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var lines []parser.LogLine
	for _, raw := range splitLines(string(data)) {
		if raw == "" {
			continue
		}
		lines = append(lines, parser.ParseLine(raw))
	}
	return lines, nil
}
