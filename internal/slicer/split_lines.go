package slicer

import "strings"

// splitLines splits a log blob into individual non-empty trimmed lines.
// It is a shared utility used by merge, dedup, and other slicers that
// operate on in-memory string content.
func splitLines(content string) []string {
	raw := strings.Split(content, "\n")
	out := make([]string, 0, len(raw))
	for _, l := range raw {
		trimmed := strings.TrimRight(l, "\r")
		out = append(out, trimmed)
	}
	return out
}
