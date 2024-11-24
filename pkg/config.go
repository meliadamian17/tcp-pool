package pool

import (
	"time"

	"github.com/meliadamian17/tcp-pool/internal"
	"github.com/meliadamian17/tcp-pool/internal/backoff"
	"github.com/meliadamian17/tcp-pool/utils"
)

// Config represents the configuration for a connection pool.
// It encapsulates various settings like connection limits, timeouts, retries,
// backoff strategies, and event hooks.
type Config struct {
	impl *internal.ConfigImpl
}

// NewConfig creates a new Config object with the specified parameters.
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
//   - A pointer to the created Config object.
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
