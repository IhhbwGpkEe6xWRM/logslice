package slicer

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"time"

	"github.com/user/logslice/internal/parser"
)

// DedupOptions configures deduplication behaviour.
type DedupOptions struct {
	// WindowSize is the number of recent line hashes to track.
	// Defaults to 1000 if zero.
	WindowSize int
}

// SliceWithDedup extracts a time-range from r, writing unique lines to w.
// Lines whose content (excluding timestamp) was seen within the rolling
// WindowSize are silently dropped.
func SliceWithDedup(r io.Reader, w io.Writer, from, to time.Time, opts DedupOptions) (int, error) {
	if opts.WindowSize <= 0 {
		opts.WindowSize = 1000
	}

	seen := make(map[string]struct{}, opts.WindowSize)
	order := make([]string, 0, opts.WindowSize)

	scanner := bufio.NewScanner(r)
	written := 0

	for scanner.Scan() {
		line := scanner.Text()
		ll, ok := parser.ParseLine(line)
		if !ok {
			continue
		}
		if !ll.InRange(from, to) {
			continue
		}

		key := hashBody(ll.Message)
		if _, dup := seen[key]; dup {
			continue
		}

		// Evict oldest entry when window is full.
		if len(order) >= opts.WindowSize {
			oldest := order[0]
			order = order[1:]
			delete(seen, oldest)
		}
		seen[key] = struct{}{}
		order = append(order, key)

		fmt.Fprintln(w, line)
		written++
	}

	return written, scanner.Err()
}

func hashBody(msg string) string {
	h := md5.Sum([]byte(msg))
	return fmt.Sprintf("%x", h)
}
