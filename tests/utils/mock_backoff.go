package utils

import "time"

type MockBackoff struct {
	calls []uint
}

func (m *MockBackoff) NextRetry(attempt uint) time.Duration {
	m.calls = append(m.calls, attempt)
	return 0
}
