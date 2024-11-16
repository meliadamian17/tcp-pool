package backoff

import "time"

type FixedBackoff struct {
	Interval uint
}

func (b *FixedBackoff) NextRetry(attempt uint) time.Duration {
	return time.Duration(b.Interval) * time.Second
}
