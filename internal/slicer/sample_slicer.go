package slicer

import (
	"bufio"
	"io"
	"math/rand"

	"github.com/logslice/logslice/internal/parser"
)

// SampleOptions controls how log lines are sampled from the input.
type SampleOptions struct {
	// Rate is the fraction of matching lines to keep (0.0 < Rate <= 1.0).
	Rate float64
	// Seed is used to initialise the random source for reproducibility.
	// A zero value uses a default seed.
	Seed int64
}

// SampleSlice reads lines from r, keeps only those whose timestamps fall
// within [from, to], and then randomly retains each line with probability
// opts.Rate. Sampled lines are written to w.
func SampleSlice(r io.Reader, w io.Writer, from, to parser.LogTimestamp, opts SampleOptions) (int, error) {
	if opts.Rate <= 0 || opts.Rate > 1.0 {
		opts.Rate = 1.0
	}

	rng := rand.New(rand.NewSource(opts.Seed))

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1<<20), 1<<20)

	written := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		ll, err := parser.ParseLine(line)
		if err != nil {
			continue
		}

		if !ll.InRange(from, to) {
			continue
		}

		if rng.Float64() > opts.Rate {
			continue
		}

		if _, err := io.WriteString(w, line+"\n"); err != nil {
			return written, err
		}
		written++
	}

	return written, scanner.Err()
}
