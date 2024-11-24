package internal

import (
	"fmt"
	"net"
	"time"
)

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
