package http

import (
	"math"
	"time"
)

const (
	backOffMin    = 100 * time.Millisecond
	backoffMax    = 10 * time.Second
	backOffFactor = 2
)

// BackOffForAttempt calculates the backoff duration for a given attempt based on exponential backoff strategy.
// It takes the attempt number as a parameter and computes the backoff duration using the formula:
//
//	backoff = backOffMin * (backOffFactor ^ attempt)
func BackOffForAttempt(attempt float64) time.Duration {
	durf := float64(backOffMin) * math.Pow(backOffFactor, attempt)
	if durf > math.MaxInt64 {
		return backoffMax
	}
	dur := time.Duration(durf)
	//keep within bounds
	if dur < backOffMin {
		return backOffMin
	}
	if dur > backoffMax {
		return backoffMax
	}
	return dur
}
