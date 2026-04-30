package parser

import (
	"testing"
	"time"
)

func mustParseRFC(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestMultilineCollectorSingleGroup(t *testing.T) {
	c := &MultilineCollector{}
	c.AddLine("2024-01-01T10:00:00Z starting service")
	c.AddLine("  at main.go:42")
	c.AddLine("  at runner.go:10")

	groups := c.Groups()
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if len(groups[0].Lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(groups[0].Lines))
	}
}

func TestMultilineCollectorMultipleGroups(t *testing.T) {
	c := &MultilineCollector{}
	c.AddLine("2024-01-01T10:00:00Z first event")
	c.AddLine("continuation of first")
	c.AddLine("2024-01-01T10:01:00Z second event")
	c.AddLine("continuation of second")

	groups := c.Groups()
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Lines[1] != "continuation of first" {
		t.Errorf("unexpected continuation line: %q", groups[0].Lines[1])
	}
}

func TestMultilineCollectorDropsLeadingContinuation(t *testing.T) {
	c := &MultilineCollector{}
	c.AddLine("orphan line without timestamp")
	c.AddLine("2024-01-01T10:00:00Z real event")

	groups := c.Groups()
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if len(groups[0].Lines) != 1 {
		t.Errorf("orphan line should have been dropped")
	}
}

func TestMultilineGroupJoinedText(t *testing.T) {
	g := &MultilineGroup{
		Timestamp: mustParseRFC("2024-01-01T10:00:00Z"),
		Lines:     []string{"line one", "line two", "line three"},
	}
	want := "line one\nline two\nline three"
	if got := g.JoinedText(); got != want {
		t.Errorf("JoinedText() = %q, want %q", got, want)
	}
}

func TestGroupsInRange(t *testing.T) {
	c := &MultilineCollector{}
	c.AddLine("2024-01-01T09:00:00Z too early")
	c.AddLine("2024-01-01T10:00:00Z in range start")
	c.AddLine("2024-01-01T10:30:00Z in range mid")
	c.AddLine("2024-01-01T11:00:00Z in range end")
	c.AddLine("2024-01-01T12:00:00Z too late")

	from := mustParseRFC("2024-01-01T10:00:00Z")
	to := mustParseRFC("2024-01-01T11:00:00Z")

	result := c.GroupsInRange(from, to)
	if len(result) != 3 {
		t.Fatalf("expected 3 groups in range, got %d", len(result))
	}
	for _, g := range result {
		if g.Timestamp.Before(from) || g.Timestamp.After(to) {
			t.Errorf("group timestamp %v out of range", g.Timestamp)
		}
	}
}

func TestGroupsInRangeEmpty(t *testing.T) {
	c := &MultilineCollector{}
	from := mustParseRFC("2024-01-01T10:00:00Z")
	to := mustParseRFC("2024-01-01T11:00:00Z")
	if got := c.GroupsInRange(from, to); len(got) != 0 {
		t.Errorf("expected empty result, got %d groups", len(got))
	}
}
