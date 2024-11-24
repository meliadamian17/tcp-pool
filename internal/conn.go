package internal

import (
	"fmt"
	"net"
	"time"
)

// Validate checks if a given connection is still valid by performing a read operation with a timeout.
// It determines whether the connection is idle, invalid, or active.
//
// Parameters:
//   - c: The connection to validate.
//
// Returns:
//   - A boolean indicating whether the connection is valid.
func Validate(c net.Conn) bool {
	c.SetReadDeadline(time.Now().Add(1 * time.Second))
	one := make([]byte, 1)
	_, err := c.Read(one)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("Connection is idle (timeout occurred)")
			return false
		}
		fmt.Printf("Connection is invalid: %v\n", err)
		return false
	}
	fmt.Println("Connection is valid (data received)")
	return true
}

// CleanupIdleConns periodically checks for idle connections and removes them if they are no longer valid.
func (p *ConnectionPool) CleanupIdleConns() {
	ticker := time.NewTicker(p.IdleTimeout)
	defer ticker.Stop()
	for range ticker.C {
		numIdle := len(p.IdleConns)
		fmt.Printf("Idle connections to process: %d\n", numIdle)

		for i := 0; i < numIdle; i++ {
			select {
			case c := <-p.IdleConns:
				if !Validate(c) {
					fmt.Println("Cleaning up idle connection")
					c.Close()
					if p.Hooks.OnConnectionClose != nil {
						p.Hooks.OnConnectionClose(c)
					}
				} else {
					fmt.Println("Connection is valid, requeuing")
					p.IdleConns <- c
				}
			default:
				fmt.Println("No more idle connections to process")
			}
		}
	}
}
