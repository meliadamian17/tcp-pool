package internal

import (
	"meliadamian17/tcp-pool/utils"
	"time"
)

type ConfigImpl struct {
	address        string
	name           string
	maxConnections int
	connTimeout    time.Duration
	idleTimeout    time.Duration
	Hooks          PoolHooks
}

func NewConfig(
	address, name string,
	maxConnections int,
	connTimeout, idleTimeout time.Duration,
	hooks PoolHooks,
) *ConfigImpl {
	if len(name) == 0 {
		name = utils.IDByAddress(address)
	}
	return &ConfigImpl{
		address:        address,
		name:           name,
		maxConnections: maxConnections,
		connTimeout:    connTimeout,
		idleTimeout:    idleTimeout,
	}
}
