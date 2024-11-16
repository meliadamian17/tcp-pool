package backoff

import "time"

type LinearBackoff struct {
	Scalar uint
}

func (b *LinearBackoff) NextRetry(attempt uint) time.Duration {
	return time.Duration(b.Scalar*attempt) * time.Second
}
