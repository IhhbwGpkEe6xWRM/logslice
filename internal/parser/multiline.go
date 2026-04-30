package parser

import (
	"strings"
	"time"
)

// MultilineGroup represents a group of log lines that belong together,
// where only the first line has a timestamp (e.g. stack traces).
type MultilineGroup struct {
	Timestamp time.Time
	Lines     []string
}

// JoinedText returns all lines joined by newline.
func (g *MultilineGroup) JoinedText() string {
	return strings.Join(g.Lines, "\n")
}

// MultilineCollector accumulates lines into groups. Lines without a
// detectable timestamp are attached to the most recently opened group.
type MultilineCollector struct {
	groups  []*MultilineGroup
	current *MultilineGroup
}

// AddLine adds a raw log line to the collector.
// If the line carries a timestamp it starts a new group; otherwise it
// is appended to the current group. Lines that arrive before any
// timestamped line is seen are silently dropped.
func (c *MultilineCollector) AddLine(line string) {
	ts, ok := TryParseTimestampPrefix(line)
	if ok {
		g := &MultilineGroup{
			Timestamp: ts,
			Lines:     []string{line},
		}
		c.groups = append(c.groups, g)
		c.current = g
		return
	}
	if c.current != nil {
		c.current.Lines = append(c.current.Lines, line)
	}
}

// Groups returns all collected multiline groups.
func (c *MultilineCollector) Groups() []*MultilineGroup {
	return c.groups
}

// GroupsInRange returns groups whose timestamp falls within [from, to].
func (c *MultilineCollector) GroupsInRange(from, to time.Time) []*MultilineGroup {
	var result []*MultilineGroup
	for _, g := range c.groups {
		if (g.Timestamp.Equal(from) || g.Timestamp.After(from)) &&
			(g.Timestamp.Equal(to) || g.Timestamp.Before(to)) {
			result = append(result, g)
		}
	}
	return result
}
