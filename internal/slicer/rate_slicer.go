package slicer

import (
	"bufio"
	"io"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// RateWindow groups log lines into fixed time buckets and counts events per window.
type RateWindow struct {
	Bucket    time.Time
	Count     int
	FirstLine string
}

// SliceRate reads log lines from r within [from, to] and aggregates them into
// fixed-duration buckets, returning one RateWindow per bucket.
func SliceRate(r io.Reader, from, to time.Time, window time.Duration) ([]RateWindow, error) {
	if window <= 0 {
		window = time.Minute
	}

	scanner := bufio.NewScanner(r)
	buckets := make(map[time.Time]*RateWindow)
	var order []time.Time

	for scanner.Scan() {
		raw := scanner.Text()
		if raw == "" {
			continue
		}

		line := parser.ParseLine(raw)
		if line.Timestamp.IsZero() {
			continue
		}
		if line.Timestamp.Before(from) || line.Timestamp.After(to) {
			continue
		}

		bucket := line.Timestamp.Truncate(window)
		if _, exists := buckets[bucket]; !exists {
			buckets[bucket] = &RateWindow{
				Bucket:    bucket,
				FirstLine: raw,
			}
			order = append(order, bucket)
		}
		buckets[bucket].Count++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	result := make([]RateWindow, 0, len(order))
	for _, b := range order {
		result = append(result, *buckets[b])
	}
	return result, nil
}
