package backoff

import (
	"math"
	"time"
)

type ExponentialBackoff struct {
	Base     uint
	MaxDelay uint
}

func (b *ExponentialBackoff) NextRetry(attempt uint) time.Duration {
	delay := time.Duration(b.Base) * time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
	if delay > time.Duration(b.MaxDelay)*time.Second {
		return time.Duration(b.MaxDelay) * time.Second
	}
	return delay
}
