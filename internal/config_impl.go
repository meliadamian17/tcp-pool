package internal

import (
	"time"

	"github.com/meliadamian17/tcp-pool/internal/backoff"
	"github.com/meliadamian17/tcp-pool/utils"
)

// ConfigImpl holds the internal configuration for the connection pool.
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

// NewConfig creates a new ConfigImpl instance.
//
// Parameters:
//   - address: The network address for the pool's connections.
//   - name: A custom name for the pool (if empty, it's generated based on the address).
//   - maxConnections: The maximum number of active connections in the pool.
//   - connTimeout: The timeout for establishing new connections.
//   - idleTimeout: The timeout for cleaning up idle connections.
//   - maxRetries: The maximum number of retries for failed connections.
//   - backoff: The backoff strategy for retrying failed connections.
//   - hooks: A set of custom hooks for pool events.
//
// Returns:
//   - A pointer to the created ConfigImpl.
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
