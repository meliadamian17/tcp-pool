package internal

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/netip"
	"sync"
	"time"
)

type ConnectionPool struct {
	address        string
	name           string
	maxConnections int
	idleTimeout    time.Duration
	connTimeout    time.Duration
	mu             sync.Mutex
	idleConns      chan net.Conn
	activeConns    int
	outputStream   *log.Logger
}

func NewConnectionPool(
	address, name string,
	maxConns int,
	connTimeout,
	idleTimeout time.Duration,
	outputStream *log.Logger,
) (*ConnectionPool, error) {

	if maxConns < 0 {
		return nil, errors.New(
			fmt.Sprintf("Max Conns Must be greater than 0. Supplied: %v", maxConns),
		)
	}

	if _, err := netip.ParseAddr(address); err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid IP Supplied. | %v", address))
	}

	pool := &ConnectionPool{
		address:        address,
		name:           name,
		maxConnections: maxConns,
		connTimeout:    connTimeout,
		idleTimeout:    idleTimeout,
		idleConns:      make(chan net.Conn, maxConns),
		outputStream:   outputStream,
	}
	go pool.CleanupIdleConns()
	return pool, nil
}

func (p *ConnectionPool) Get() (net.Conn, error) {
	select {
	case conn := <-p.idleConns:
		// Retrieved an idle connection, check if it's still valid.
		if Validate(conn) {
			p.outputStream.Printf("Idle Connnection Found to %v", p.address)
			return conn, nil
		}
		// If not valid, try creating a new connection
		p.outputStream.Printf("No Idle Connection Found! Trying to Open a New Connection...")
		return p.newConnection()
	default:
		// No idle connections, try to create a new one if under max limit.
		p.outputStream.Printf("No Idle Connection Found! Trying to Open a New Connection...")
		return p.newConnection()
	}
}

// TODO: Add retry options (either retry spawning connection or go back and look for an idle connection again)
func (p *ConnectionPool) newConnection() (net.Conn, error) {
	conn, err := net.Dial("tcp", p.address)
	if err != nil {
		p.outputStream.Printf("Unable to create a new connection to %v", p.address)
		return nil, err
	}

	return conn, nil
}

func (p *ConnectionPool) Release(conn net.Conn) error {
	select {
	case p.idleConns <- conn:
		// Successfully returned the connection to the pool.
		p.outputStream.Print("Successfully released connection back into the pool")
		return nil
	default:
		// Pool is full, close the connection.
		p.outputStream.Printf("The pool is full! Connection was is now closed")
		return conn.Close()
	}
}
