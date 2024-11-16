package pool

import "meliadamian17/tcp-pool/internal/backoff"

func NewExponentialBackoff(baseDelay, maxDelay uint) backoff.ExponentialBackoff {
	return backoff.ExponentialBackoff{
		Base:     baseDelay,
		MaxDelay: maxDelay,
	}
}

func NewFibonacciBackoff(maxDelay uint) backoff.FibonacciBackoff {
	return backoff.FibonacciBackoff{
		MaxDelay: maxDelay,
	}
}

func NewFixedBackoff(interval uint) backoff.FixedBackoff {
	return backoff.FixedBackoff{
		Interval: interval,
	}
}

func NewLinearBackoff(scalar uint) backoff.LinearBackoff {
	return backoff.LinearBackoff{
		Scalar: scalar,
	}
}

func NewPolynomialBackoff(exponent uint) backoff.PolynomialBackoff {
	return backoff.PolynomialBackoff{
		Exponent: exponent,
	}
}
