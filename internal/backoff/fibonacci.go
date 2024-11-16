package backoff

import "time"

var FibonacciMap = map[uint]uint{
	0: 0,
	1: 1,
	2: 1,
	3: 2,
	4: 3,
	5: 5,
	6: 8,
	7: 13,
	8: 21,
	9: 34,
}

type FibonacciBackoff struct {
	MaxDelay uint
}

func (b *FibonacciBackoff) NextRetry(attempt uint) time.Duration {
	delay := time.Duration(fib(attempt, FibonacciMap)) * time.Second
	if delay > time.Duration(b.MaxDelay) {
		return time.Duration(b.MaxDelay) * time.Second
	}

	return time.Duration(delay)
}

func fib(n uint, memo map[uint]uint) uint {
	if val, ok := memo[n]; ok {
		return val
	}

	if n == 0 || n == 1 {
		return n
	}

	memo[n] = fib(n-1, memo) + fib(n-2, memo)
	return memo[n]
}
