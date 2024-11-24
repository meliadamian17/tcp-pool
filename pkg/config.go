package pool

import (
	"meliadamian17/tcp-pool/internal"
	"meliadamian17/tcp-pool/internal/backoff"
	"meliadamian17/tcp-pool/utils"
	"time"
)

type Config struct {
	impl *internal.ConfigImpl
}

func NewConfig(
	address, name string,
	maxConnections int,
	connTimeout, idleTimeout time.Duration,
	maxRetries uint,
	backoff backoff.Backoff,
	hooks PoolHooks,
) *Config {
	if len(name) == 0 {
		name = utils.IDByAddress(address)
	}
	impl := internal.NewConfig(
		address,
		name,
		maxConnections,
		connTimeout,
		idleTimeout,
		maxRetries,
		backoff,
		hooks.ToInternal(),
	)
	return &Config{impl: impl}
}
