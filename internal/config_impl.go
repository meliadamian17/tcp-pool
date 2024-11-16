package internal

import (
	"meliadamian17/tcp-pool/internal/backoff"
	"meliadamian17/tcp-pool/utils"
	"time"
)

type ConfigImpl struct {
	Address        string
	Name           string
	MaxConnections int
	ConnTimeout    time.Duration
	IdleTimeout    time.Duration
	MaxRetries     uint
	Backoff        backoff.Backoff
	Hooks          PoolHooks
}

func NewConfig(
	address, name string,
	maxConnections int,
	connTimeout, idleTimeout time.Duration,
	maxRetries uint,
	backoff backoff.Backoff,
	hooks PoolHooks,
) *ConfigImpl {
	if len(name) == 0 {
		name = utils.IDByAddress(address)
	}
	return &ConfigImpl{
		Address:        address,
		Name:           name,
		MaxConnections: maxConnections,
		ConnTimeout:    connTimeout,
		IdleTimeout:    idleTimeout,
		MaxRetries:     maxRetries,
		Backoff:        backoff,
		Hooks:          hooks,
	}
}
