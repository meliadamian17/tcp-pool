package pool

import (
	"meliadamian17/tcp-pool/internal"
	"net"
)

type Pool struct {
	impl *internal.ConnectionPool
}

func New(c Config) (*Pool, error) {
	impl, err := internal.NewConnectionPool(*c.impl)
	if err != nil {
		return nil, err
	}
	return &Pool{impl: impl}, nil
}

func (p *Pool) Get() (net.Conn, error) {
	return p.impl.Get()
}

func (p *Pool) GetAsync() <-chan struct {
	Conn net.Conn
	Err  error
} {
	resultChan := make(chan struct {
		Conn net.Conn
		Err  error
	})

	go func() {
		conn, err := p.impl.Get()
		resultChan <- struct {
			Conn net.Conn
			Err  error
		}{Conn: conn, Err: err}
		close(resultChan)
	}()

	return resultChan
}

// Release returns a connection to the pool.
func (p *Pool) Release(conn net.Conn) error {
	return p.impl.Release(conn)
}
