package http

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestBackOffForAttempt(t *testing.T) {
	tests := []struct {
		attempt float64
		result  time.Duration
	}{
		{0, backOffMin},
		{1, backOffMin * time.Duration(backOffFactor)},
		{2, backOffMin * time.Duration(math.Pow(backOffFactor, 2))},
		{3, backOffMin * time.Duration(math.Pow(backOffFactor, 3))},
		{10, backoffMax},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Attempt %v", test.attempt), func(t *testing.T) {
			result := BackOffForAttempt(test.attempt)

			const delta = time.Millisecond
			if result < test.result-delta || result > test.result+delta {
				t.Errorf("expected %v, got %v", test.result, result)
			}

			if result < backOffMin {
				t.Errorf("result %v is less than backOffMin %v", result, backOffMin)
			}
			if result > backoffMax {
				t.Errorf("result %v is greater than backoffMax %v", result, backoffMax)
			}
		})
	}
}
