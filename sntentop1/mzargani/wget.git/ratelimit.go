package main

import (
	"io"
	"time"
)

// rateLimitedReader wraps an io.Reader and enforces a bytes-per-second limit.
type rateLimitedReader struct {
	r         io.Reader
	rateLimit int64 // bytes per second
	lastRead  time.Time
	readSoFar int64
	startTime time.Time
}

func (r *rateLimitedReader) Read(p []byte) (int, error) {
	if r.startTime.IsZero() {
		r.startTime = time.Now()
	}

	// Limit chunk size to avoid reading too much at once
	maxChunk := r.rateLimit / 10
	if maxChunk < 1024 {
		maxChunk = 1024
	}
	if int64(len(p)) > maxChunk {
		p = p[:maxChunk]
	}

	n, err := r.r.Read(p)
	if n > 0 {
		r.readSoFar += int64(n)

		// Calculate how long we should have taken to read this much data
		expectedDuration := time.Duration(float64(r.readSoFar) / float64(r.rateLimit) * float64(time.Second))
		actualDuration := time.Since(r.startTime)

		if expectedDuration > actualDuration {
			time.Sleep(expectedDuration - actualDuration)
		}
	}

	return n, err
}
