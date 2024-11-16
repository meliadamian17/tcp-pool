package internal

import (
	"fmt"
	"net"
	"time"
)

func Validate(c net.Conn) bool {
	one := []byte{}
	c.SetReadDeadline(time.Now().Add(1 * time.Second))
	if _, err := c.Read(one); err != nil {
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
		case c := <-p.idleConns:
			if !Validate(c) {
				c.Close()
				if p.hooks.OnConnectionClose != nil {
					p.hooks.OnConnectionClose(c)
				} else {
					fmt.Printf("Connection was cleaned up due to being idle")
				}
			} else {
				// Put it back if still valid
				p.idleConns <- c
			}
		default:
			fmt.Print("There are no idle connections")
			// No connections to clean up.
		}
	}
}
