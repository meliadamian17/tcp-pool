package pool

import (
	"meliadamian17/tcp-pool/internal"
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
	hooks internal.PoolHooks,
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
		hooks,
	)
	return &Config{impl: impl}
}
