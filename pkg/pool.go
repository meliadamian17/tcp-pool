package pool

import (
	"meliadamian17/tcp-pool/internal"
	"net"
)

type Pool struct {
	impl *internal.ConnectionPool
}

func New(c Config) (*Pool, error) {
	impl, err := internal.NewConnectionPool(
		c.address,
		c.name,
		c.maxConnections,
		c.connTimeout,
		c.idleTimeout,
		c.outputStream,
	)
	if err != nil {
		return nil, err
	}
	return &Pool{impl: impl}, nil
}

func (p *Pool) Get() (net.Conn, error) {
	return p.impl.Get()
}

func (p *Pool) Release(conn net.Conn) error {
	return p.impl.Release(conn)
}
