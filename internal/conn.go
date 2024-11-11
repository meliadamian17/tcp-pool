package internal

import (
	"net"
	"time"
)

// TODO: finish using this everywhere
type ConnWithID struct {
	Conn net.Conn
	ID   string
}

func Validate(conn net.Conn) bool {
	one := []byte{}
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	if _, err := conn.Read(one); err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return true
		}
		return false
	}
	return true
}

func (p *ConnectionPool) CleanupIdleConns() {
	ticker := time.NewTicker(p.idleTimeout)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case conn := <-p.idleConns:
			if !Validate(conn) {
				conn.Close()
				//TODO: attach IDs to connections
				p.outputStream.Print("Connection was cleaned up due to being idle")
			} else {
				// Put it back if still valid
				p.idleConns <- conn
			}
		default:
			p.outputStream.Print("There are no idle connections")
			// No connections to clean up.
		}
	}
}
