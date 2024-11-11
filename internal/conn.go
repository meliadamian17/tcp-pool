package internal

import (
	"net"
	"time"
)

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
			} else {
				// Put it back if still valid
				p.idleConns <- conn
			}
		default:
			// No connections to clean up.
		}
	}
}
