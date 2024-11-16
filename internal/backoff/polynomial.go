package backoff

import (
	"math"
	"time"
)

type PolynomialBackoff struct {
	Exponent uint
}

func (b *PolynomialBackoff) NextRetry(attempt uint) time.Duration {
	return time.Duration(math.Pow(float64(attempt), float64(b.Exponent))) * time.Second
}
