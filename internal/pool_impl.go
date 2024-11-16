package internal

import (
	"errors"
	"fmt"
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
	hooks          PoolHooks
}

func NewConnectionPool(
	c ConfigImpl,
) (*ConnectionPool, error) {

	if c.maxConnections < 0 {
		if c.Hooks.OnPoolCreateError != nil {
			c.Hooks.OnPoolCreateError(errors.New(
				fmt.Sprintf("Max Conns Must be greater than 0. Supplied: %v", c.maxConnections),
			),
			)
		} else {
			fmt.Printf("Max Conns Must be greater than 0. Supplied: %v", c.maxConnections)
		}
		return nil, errors.New(
			fmt.Sprintf("Max Conns Must be greater than 0. Supplied: %v", c.maxConnections),
		)
	}

	if _, err := netip.ParseAddr(c.address); err != nil {
		if c.Hooks.OnPoolCreateError != nil {
			c.Hooks.OnPoolCreateError(errors.New(
				fmt.Sprintf("Invalid IP Supplied. | %v", c.address),
			),
			)
		} else {
			fmt.Printf(fmt.Sprintf("Invalid IP Supplied. | %v", c.address))
		}
		return nil, errors.New(fmt.Sprintf("Invalid IP Supplied. | %v", c.address))
	}

	pool := &ConnectionPool{
		address:        c.address,
		name:           c.name,
		maxConnections: c.maxConnections,
		connTimeout:    c.connTimeout,
		idleTimeout:    c.idleTimeout,
		idleConns:      make(chan net.Conn),
		hooks:          c.Hooks,
	}

	if pool.hooks.OnPoolCreate != nil {
		pool.hooks.OnPoolCreate(c)
	} else {
		fmt.Printf("New Connection pool for address %v created", pool.address)
	}

	go pool.CleanupIdleConns()

	return pool, nil
}

func (p *ConnectionPool) Get() (net.Conn, error) {
	select {
	case conn := <-p.idleConns:
		// Retrieved an idle connection, check if it's still valid.
		if Validate(conn) {
			if p.hooks.OnConnectionAcquire != nil {
				p.hooks.OnConnectionAcquire(conn)
			} else {
				fmt.Printf("Idle Connnection Found to %v", p.address)
			}
			return conn, nil
		}
		// If not valid, try creating a new connection
		fmt.Printf("No Idle Connection Found! Trying to Open a New Connection...")
		conn, err := p.newConnection()
		if err != nil {
			if p.hooks.OnConnectionError != nil {
				p.hooks.OnConnectionError(err)
			} else {
				fmt.Printf("Failed to create new connection: %v", err)
				return nil, err
			}
		}
		if p.hooks.OnConnectionCreate != nil {
			p.hooks.OnConnectionCreate(conn)
		} else {
			fmt.Printf("New Connection Created: %v", conn)
		}

		return conn, nil

	default:
		fmt.Printf("No Idle Connection Found! Trying to Open a New Connection...")
		conn, err := p.newConnection()
		if err != nil {
			if p.hooks.OnConnectionError != nil {
				p.hooks.OnConnectionError(err)
			} else {
				fmt.Printf("Failed to create new connection: %v", err)
				return nil, err
			}
		}
		if p.hooks.OnConnectionCreate != nil {
			p.hooks.OnConnectionCreate(conn)
		} else {
			fmt.Printf("New Connection Created: %v", conn)
		}

		return conn, nil
	}
}

// TODO: Add retry options (either retry spawning connection or go back and look for an idle connection again)
func (p *ConnectionPool) newConnection() (net.Conn, error) {
	conn, err := net.Dial("tcp", p.address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (p *ConnectionPool) Release(conn net.Conn) error {
	select {
	case p.idleConns <- conn:
		// Successfully returned the connection to the pool.
		if p.hooks.OnConnectionRelease != nil {
			p.hooks.OnConnectionRelease(conn)
		} else {
			fmt.Print("Successfully released connection back into the pool")
		}
		return nil
	default:
		// Pool is full, close the connection.
		if p.hooks.OnConnectionClose != nil {
			p.hooks.OnConnectionClose(conn)
		} else {
			fmt.Printf("Connection Is Closing Due to Pool being full")
		}
		return conn.Close()
	}
}
