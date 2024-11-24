package internal

import (
	"errors"
	"fmt"
	"meliadamian17/tcp-pool/internal/backoff"
	"net"
	"time"
)

// ConnectionPool represents a pool of reusable TCP connections.
// It manages the creation, reuse, and cleanup of idle connections.
type ConnectionPool struct {
	Address        string          // Network address for the pool's connections
	Name           string          // Name of the connection pool
	MaxConnections int             // Maximum number of active connections
	IdleTimeout    time.Duration   // Duration after which idle connections are cleaned up
	ConnTimeout    time.Duration   // Timeout for establishing a new connection
	IdleConns      chan net.Conn   // Channel for storing idle connections
	ActiveConns    int             // Current number of active connections
	MaxRetries     uint            // Maximum number of retries for connection establishment
	Backoff        backoff.Backoff // Backoff strategy for retries
	Hooks          PoolHooks       // Hooks for connection pool events
}

// NewConnectionPool initializes a new ConnectionPool with the given configuration.
//
// Parameters:
//   - c: A ConfigImpl object containing pool configuration.
//
// Returns:
//   - A pointer to the created ConnectionPool.
//   - An error, if the initialization fails.
func NewConnectionPool(c ConfigImpl) (*ConnectionPool, error) {
	if c.MaxConnections < 0 {
		err := errors.New(
			fmt.Sprintf("Max Conns must be greater than 0. Supplied: %v", c.MaxConnections),
		)
		if c.Hooks.OnPoolCreateError != nil {
			c.Hooks.OnPoolCreateError(err)
		} else {
			fmt.Println(err)
		}
		return nil, err
	}

	pool := &ConnectionPool{
		Address:        c.Address,
		Name:           c.Name,
		MaxConnections: c.MaxConnections,
		ConnTimeout:    c.ConnTimeout,
		IdleTimeout:    c.IdleTimeout,
		IdleConns:      make(chan net.Conn, c.MaxConnections),
		MaxRetries:     c.MaxRetries,
		Backoff:        c.Backoff,
		Hooks:          c.Hooks,
	}

	if pool.Hooks.OnPoolCreate != nil {
		pool.Hooks.OnPoolCreate(c)
	} else {
		fmt.Printf("New connection pool for address %v created\n", pool.Address)
	}

	go pool.CleanupIdleConns()

	return pool, nil
}

// Get retrieves a connection from the pool. If an idle connection is available,
// it is returned; otherwise, a new connection is created.
//
// Returns:
//   - A net.Conn object representing the connection.
//   - An error, if the connection retrieval fails.
func (p *ConnectionPool) Get() (net.Conn, error) {
	select {
	case conn := <-p.IdleConns:
		if Validate(conn) {
			if p.Hooks.OnConnectionAcquire != nil {
				p.Hooks.OnConnectionAcquire(conn)
			} else {
				fmt.Printf("Idle connection found to %v\n", p.Address)
			}
			return conn, nil
		}
		fmt.Printf("No valid idle connection found! Trying to open a new connection...\n")
		return p.newConnection()
	default:
		fmt.Printf("No idle connection found! Trying to open a new connection...\n")
		return p.newConnection()
	}
}

// newConnectionAsync creates a new connection asynchronously and applies backoff strategies for retries.
//
// Returns:
//   - A channel that sends a struct containing the connection or an error.
func (p *ConnectionPool) newConnectionAsync() <-chan struct {
	conn net.Conn
	err  error
} {
	resultChan := make(chan struct {
		conn net.Conn
		err  error
	})

	go func() {
		var conn net.Conn
		var err error

		for attempt := 1; attempt <= int(p.MaxRetries); attempt++ {
			conn, err = net.DialTimeout("tcp", p.Address, p.ConnTimeout)
			if err == nil {
				resultChan <- struct {
					conn net.Conn
					err  error
				}{conn: conn, err: nil}
				close(resultChan)
				return
			}

			time.Sleep(p.Backoff.NextRetry(uint(attempt)))
		}

		resultChan <- struct {
			conn net.Conn
			err  error
		}{conn: nil, err: fmt.Errorf("failed to establish connection after %d retries: %w", p.MaxRetries, err)}
		close(resultChan)
	}()

	return resultChan
}

// newConnection creates a new connection synchronously and triggers hooks for connection events.
//
// Returns:
//   - A net.Conn object representing the connection.
//   - An error, if the connection creation fails.
func (p *ConnectionPool) newConnection() (net.Conn, error) {
	resultChan := p.newConnectionAsync()
	result := <-resultChan
	if result.err != nil {
		if p.Hooks.OnConnectionError != nil {
			p.Hooks.OnConnectionError(result.err)
		} else {
			fmt.Printf("Failed to create new connection: %v\n", result.err)
		}
		return nil, result.err
	}

	if p.Hooks.OnConnectionCreate != nil {
		p.Hooks.OnConnectionCreate(result.conn)
	} else {
		fmt.Printf("New connection created: %v\n", result.conn)
	}

	return result.conn, nil
}

// Release returns a previously acquired connection to the pool.
// If the pool is full, the connection is closed instead.
//
// Parameters:
//   - conn: The connection to be returned to the pool.
//
// Returns:
//   - An error, if the release process fails.
func (p *ConnectionPool) Release(conn net.Conn) error {
	select {
	case p.IdleConns <- conn:
		if p.Hooks.OnConnectionRelease != nil {
			p.Hooks.OnConnectionRelease(conn)
		} else {
			fmt.Println("Successfully released connection back into the pool")
		}
		return nil
	default:
		if p.Hooks.OnConnectionClose != nil {
			p.Hooks.OnConnectionClose(conn)
		} else {
			fmt.Println("Connection is closing due to pool being full")
		}
		return conn.Close()
	}
}
