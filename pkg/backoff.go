package pool

import "meliadamian17/tcp-pool/internal/backoff"

func NewExponentialBackoff(baseDelay, maxDelay uint) backoff.Backoff {
	return &backoff.ExponentialBackoff{
		Base:     baseDelay,
		MaxDelay: maxDelay,
	}
}

func NewFibonacciBackoff(maxDelay uint) backoff.Backoff {
	return &backoff.FibonacciBackoff{
		MaxDelay: maxDelay,
	}
}

func NewFixedBackoff(interval uint) backoff.Backoff {
	return &backoff.FixedBackoff{
		Interval: interval,
	}
}

func NewLinearBackoff(scalar uint) backoff.Backoff {
	return &backoff.LinearBackoff{
		Scalar: scalar,
	}
}

func NewPolynomialBackoff(exponent uint) backoff.Backoff {
	return &backoff.PolynomialBackoff{
		Exponent: exponent,
	}
}
