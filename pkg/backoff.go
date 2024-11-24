package pool

import "meliadamian17/tcp-pool/internal/backoff"

// NewExponentialBackoff creates a new exponential backoff strategy.
// The delay between retries doubles with each attempt until reaching the maximum delay.
//
// Parameters:
//   - baseDelay: The initial delay in seconds for the first retry.
//   - maxDelay: The maximum delay in seconds for retries.
//
// Returns:
//   - A backoff.Backoff implementation using exponential backoff.
func NewExponentialBackoff(baseDelay, maxDelay uint) backoff.Backoff {
	return &backoff.ExponentialBackoff{
		Base:     baseDelay,
		MaxDelay: maxDelay,
	}
}

// NewFibonacciBackoff creates a new Fibonacci backoff strategy.
// The delay between retries follows the Fibonacci sequence until reaching the maximum delay.
//
// Parameters:
//   - maxDelay: The maximum delay in seconds for retries.
//
// Returns:
//   - A backoff.Backoff implementation using Fibonacci backoff.
func NewFibonacciBackoff(maxDelay uint) backoff.Backoff {
	return &backoff.FibonacciBackoff{
		MaxDelay: maxDelay,
	}
}

// NewFixedBackoff creates a new fixed backoff strategy.
// The delay between retries remains constant.
//
// Parameters:
//   - interval: The fixed delay in seconds between retries.
//
// Returns:
//   - A backoff.Backoff implementation using fixed backoff.
func NewFixedBackoff(interval uint) backoff.Backoff {
	return &backoff.FixedBackoff{
		Interval: interval,
	}
}

// NewLinearBackoff creates a new linear backoff strategy.
// The delay between retries increases linearly with each attempt.
//
// Parameters:
//   - scalar: The constant value added to the delay for each retry.
//
// Returns:
//   - A backoff.Backoff implementation using linear backoff.
func NewLinearBackoff(scalar uint) backoff.Backoff {
	return &backoff.LinearBackoff{
		Scalar: scalar,
	}
}

// NewPolynomialBackoff creates a new polynomial backoff strategy.
// The delay between retries follows a polynomial growth pattern.
//
// Parameters:
//   - exponent: The exponent used for calculating delay growth.
//
// Returns:
//   - A backoff.Backoff implementation using polynomial backoff.
func NewPolynomialBackoff(exponent uint) backoff.Backoff {
	return &backoff.PolynomialBackoff{
		Exponent: exponent,
	}
}
