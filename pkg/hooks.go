package pool

import (
	"meliadamian17/tcp-pool/internal"
	"net"
)

// PoolHooks defines callback functions that can be triggered
// during various connection pool events. These hooks allow custom
// logic to be executed during connection creation, acquisition, release,
// closure, errors, or during pool creation.
type PoolHooks struct {
	// OnConnectionCreate is triggered when a new connection is created.
	OnConnectionCreate func(conn net.Conn)
	// OnConnectionAcquire is triggered when a connection is acquired from the pool.
	OnConnectionAcquire func(conn net.Conn)
	// OnConnectionRelease is triggered when a connection is released back into the pool.
	OnConnectionRelease func(conn net.Conn)
	// OnConnectionClose is triggered when a connection is closed.
	OnConnectionClose func(conn net.Conn)
	// OnConnectionError is triggered when an error occurs during connection operations.
	OnConnectionError func(err error)
	// OnPoolCreate is triggered when the connection pool is successfully created.
	OnPoolCreate func(c Config)
	// OnPoolCreateError is triggered when there is an error during pool creation.
	OnPoolCreateError func(err error)
}

// ToInternal converts a public PoolHooks object to the corresponding internal representation.
// This is used to pass hooks from the public API to the internal pool implementation.
func (h PoolHooks) ToInternal() internal.PoolHooks {
	return internal.PoolHooks{
		OnConnectionCreate:  h.OnConnectionCreate,
		OnConnectionAcquire: h.OnConnectionAcquire,
		OnConnectionRelease: h.OnConnectionRelease,
		OnConnectionClose:   h.OnConnectionClose,
		OnConnectionError:   h.OnConnectionError,
		OnPoolCreate: func(c internal.ConfigImpl) {
			if h.OnPoolCreate != nil {
				h.OnPoolCreate(Config{
					impl: &c,
				})
			}
		},
		OnPoolCreateError: h.OnPoolCreateError,
	}
}
