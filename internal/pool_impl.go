package internal

import (
	"errors"
	"fmt"
	"meliadamian17/tcp-pool/internal/backoff"
	"net"
	"net/netip"
	"time"
)

type ConnectionPool struct {
	address        string
	name           string
	maxConnections int
	idleTimeout    time.Duration
	connTimeout    time.Duration
	idleConns      chan net.Conn
	activeConns    int
	maxRetries     uint
	backoff        backoff.Backoff
	hooks          PoolHooks
}

func NewConnectionPool(c ConfigImpl) (*ConnectionPool, error) {
	if c.MaxConnections < 0 {
		err := errors.New(
			fmt.Sprintf("Max Conns Must be greater than 0. Supplied: %v", c.MaxConnections),
		)
		if c.Hooks.OnPoolCreateError != nil {
			c.Hooks.OnPoolCreateError(err)
		} else {
			fmt.Println(err)
		}
		return nil, err
	}

	if _, err := netip.ParseAddr(c.Address); err != nil {
		err := errors.New(fmt.Sprintf("Invalid IP Supplied. | %v", c.Address))
		if c.Hooks.OnPoolCreateError != nil {
			c.Hooks.OnPoolCreateError(err)
		} else {
			fmt.Println(err)
		}
		return nil, err
	}

	pool := &ConnectionPool{
		address:        c.Address,
		name:           c.Name,
		maxConnections: c.MaxConnections,
		connTimeout:    c.ConnTimeout,
		idleTimeout:    c.IdleTimeout,
		idleConns:      make(chan net.Conn, c.MaxConnections),
		maxRetries:     c.MaxRetries,
		backoff:        c.Backoff,
		hooks:          c.Hooks,
	}

	if pool.hooks.OnPoolCreate != nil {
		pool.hooks.OnPoolCreate(c)
	} else {
		fmt.Printf("New Connection pool for address %v created\n", pool.address)
	}

	go pool.CleanupIdleConns()

	return pool, nil
}

func (p *ConnectionPool) Get() (net.Conn, error) {
	select {
	case conn := <-p.idleConns:
		// Check if the idle connection is still valid
		if Validate(conn) {
			if p.hooks.OnConnectionAcquire != nil {
				p.hooks.OnConnectionAcquire(conn)
			} else {
				fmt.Printf("Idle Connection Found to %v\n", p.address)
			}
			return conn, nil
		}
		// Connection is invalid, create a new one
		fmt.Printf("No valid Idle Connection Found! Trying to Open a New Connection...\n")
		return p.newConnection()

	default:
		// No idle connection available, create a new one
		fmt.Printf("No Idle Connection Found! Trying to Open a New Connection...\n")
		return p.newConnection()
	}
}

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

		for attempt := 1; attempt <= int(p.maxRetries); attempt++ {
			conn, err = net.DialTimeout("tcp", p.address, p.connTimeout)
			if err == nil {
				// Successfully established the connection
				resultChan <- struct {
					conn net.Conn
					err  error
				}{conn: conn, err: nil}
				close(resultChan)
				return
			}

			// Apply backoff strategy
			time.Sleep(p.backoff.NextRetry(uint(attempt)))
		}

		// After exhausting retries, send the error
		resultChan <- struct {
			conn net.Conn
			err  error
		}{conn: nil, err: fmt.Errorf("failed to establish connection after %d retries: %w", p.maxRetries, err)}
		close(resultChan)
	}()

	return resultChan
}

func (p *ConnectionPool) newConnection() (net.Conn, error) {
	// Wait for the result from the async connection creator
	resultChan := p.newConnectionAsync()
	result := <-resultChan
	if result.err != nil {
		if p.hooks.OnConnectionError != nil {
			p.hooks.OnConnectionError(result.err)
		} else {
			fmt.Printf("Failed to create new connection: %v\n", result.err)
		}
		return nil, result.err
	}

	if p.hooks.OnConnectionCreate != nil {
		p.hooks.OnConnectionCreate(result.conn)
	} else {
		fmt.Printf("New Connection Created: %v\n", result.conn)
	}

	return result.conn, nil
}

func (p *ConnectionPool) Release(conn net.Conn) error {
	select {
	case p.idleConns <- conn:
		// Successfully returned the connection to the pool
		if p.hooks.OnConnectionRelease != nil {
			p.hooks.OnConnectionRelease(conn)
		} else {
			fmt.Println("Successfully released connection back into the pool")
		}
		return nil
	default:
		// Pool is full, close the connection
		if p.hooks.OnConnectionClose != nil {
			p.hooks.OnConnectionClose(conn)
		} else {
			fmt.Println("Connection is closing due to pool being full")
		}
		return conn.Close()
	}
}
