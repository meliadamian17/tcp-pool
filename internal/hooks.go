package internal

import (
	"net"
)

type PoolHooks struct {
	OnConnectionCreate  func(conn net.Conn)
	OnConnectionAcquire func(conn net.Conn)
	OnConnectionRelease func(conn net.Conn)
	OnConnectionClose   func(conn net.Conn)
	OnConnectionError   func(err error)
	OnPoolCreate        func(c ConfigImpl)
	OnPoolCreateError   func(err error)
}
