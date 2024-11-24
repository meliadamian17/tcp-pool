package pool

import (
	"meliadamian17/tcp-pool/internal"
	"net"
)

type PoolHooks struct {
	OnConnectionCreate  func(conn net.Conn)
	OnConnectionAcquire func(conn net.Conn)
	OnConnectionRelease func(conn net.Conn)
	OnConnectionClose   func(conn net.Conn)
	OnConnectionError   func(err error)
	OnPoolCreate        func(c Config)
	OnPoolCreateError   func(err error)
}

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
