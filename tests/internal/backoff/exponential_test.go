package internal

import (
	"testing"
	"time"

	"meliadamian17/tcp-pool/internal/backoff"
	"meliadamian17/tcp-pool/tests/utils"
)

func TestExponentialBackoff(t *testing.T) {
	b := &backoff.ExponentialBackoff{Base: 2, MaxDelay: 10}

	utils.AssertEqual(
		t,
		time.Duration(2)*time.Second,
		b.NextRetry(1),
		"First retry duration mismatch",
	)
	utils.AssertEqual(
		t,
		time.Duration(4)*time.Second,
		b.NextRetry(2),
		"Second retry duration mismatch",
	)
	utils.AssertEqual(
		t,
		time.Duration(8)*time.Second,
		b.NextRetry(3),
		"Third retry duration mismatch",
	)
	utils.AssertEqual(
		t,
		time.Duration(10)*time.Second,
		b.NextRetry(4),
		"Retry duration should be capped at MaxDelay",
	)
}
