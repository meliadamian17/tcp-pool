package internal

import (
	"errors"
	"fmt"
	"meliadamian17/tcp-pool/internal/backoff"
	"net"
	"time"
)

type ConnectionPool struct {
	Address        string
	Name           string
	MaxConnections int
	IdleTimeout    time.Duration
	ConnTimeout    time.Duration
	IdleConns      chan net.Conn
	ActiveConns    int
	MaxRetries     uint
	Backoff        backoff.Backoff
	Hooks          PoolHooks
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
		fmt.Printf("New Connection pool for address %v created\n", pool.Address)
	}

	go pool.CleanupIdleConns()

	return pool, nil
}

func (p *ConnectionPool) Get() (net.Conn, error) {
	select {
	case conn := <-p.IdleConns:
		// Check if the idle connection is still valid
		if Validate(conn) {
			if p.Hooks.OnConnectionAcquire != nil {
				p.Hooks.OnConnectionAcquire(conn)
			} else {
				fmt.Printf("Idle Connection Found to %v\n", p.Address)
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

		for attempt := 1; attempt <= int(p.MaxRetries); attempt++ {
			conn, err = net.DialTimeout("tcp", p.Address, p.ConnTimeout)
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
			time.Sleep(p.Backoff.NextRetry(uint(attempt)))
		}

		// After exhausting retries, send the error
		resultChan <- struct {
			conn net.Conn
			err  error
		}{conn: nil, err: fmt.Errorf("failed to establish connection after %d retries: %w", p.MaxRetries, err)}
		close(resultChan)
	}()

	return resultChan
}

func (p *ConnectionPool) newConnection() (net.Conn, error) {
	// Wait for the result from the async connection creator
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
		fmt.Printf("New Connection Created: %v\n", result.conn)
	}

	return result.conn, nil
}

func (p *ConnectionPool) Release(conn net.Conn) error {
	select {
	case p.IdleConns <- conn:
		// Successfully returned the connection to the pool
		if p.Hooks.OnConnectionRelease != nil {
			p.Hooks.OnConnectionRelease(conn)
		} else {
			fmt.Println("Successfully released connection back into the pool")
		}
		return nil
	default:
		// Pool is full, close the connection
		if p.Hooks.OnConnectionClose != nil {
			p.Hooks.OnConnectionClose(conn)
		} else {
			fmt.Println("Connection is closing due to pool being full")
		}
		return conn.Close()
	}
}
