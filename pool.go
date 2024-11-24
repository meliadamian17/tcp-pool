package tcppool

import (
	"net"

	"github.com/meliadamian17/tcp-pool/internal"
)

// Pool represents a connection pool that manages TCP connections.
// It provides methods to acquire and release connections, as well as to fetch them asynchronously.
type Pool struct {
	impl *internal.ConnectionPool
}

// New creates a new Pool instance based on the given configuration.
// It initializes an internal connection pool and returns a Pool object.
//
// Parameters:
//   - c: A Config object containing the pool configuration.
//
// Returns:
//   - A pointer to the created Pool.
//   - An error, if the pool initialization fails.
func New(c Config) (*Pool, error) {
	impl, err := internal.NewConnectionPool(*c.impl)
	if err != nil {
		return nil, err
	}
	return &Pool{impl: impl}, nil
}

// Get retrieves a connection from the pool.
// If an idle connection is available, it is returned; otherwise, a new connection is created.
//
// Returns:
//   - A net.Conn representing the connection.
//   - An error, if the connection retrieval fails.
func (p *Pool) Get() (net.Conn, error) {
	return p.impl.Get()
}

// GetAsync retrieves a connection from the pool asynchronously.
// It returns a channel through which the result (connection or error) will be sent once available.
//
// Returns:
//   - A read-only channel of a struct containing:
//   - Conn: A net.Conn representing the connection.
//   - Err: An error, if the connection retrieval fails.
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

// Release returns a previously acquired connection back to the pool.
// If the pool is full, the connection is closed.
//
// Parameters:
//   - conn: The connection to be released.
//
// Returns:
//   - An error, if the release process fails.
func (p *Pool) Release(conn net.Conn) error {
	return p.impl.Release(conn)
}
