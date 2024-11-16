package backoff

import "time"

type Backoff interface {
	NextRetry(attempt uint) time.Duration
}
