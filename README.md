# tcp-pool

[![Go Reference](https://pkg.go.dev/badge/github.com/meliadamian17/tcppool.svg)](https://pkg.go.dev/github.com/meliadamian17/tcppool)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`tcppool` is a lightweight and flexible Go library for managing TCP connection pooling. It simplifies connection reuse, reduces resource usage, and provides customizable hooks for various connection lifecycle events.

---

## Features

- **Connection Pooling**: Manage TCP connections efficiently.
- **Customizable Backoff Strategies**: Includes exponential, Fibonacci, linear, polynomial, and fixed backoff.
- **Lifecycle Hooks**: Add custom logic for connection creation, acquisition, release, and errors.
- **Idle Connection Cleanup**: Automatically removes stale or invalid connections.
- **Asynchronous Connection Retrieval**: Fetch connections asynchronously when needed.

---

## Installation

Install the library using `go get`:

```bash
go get github.com/meliadamian17/tcp-pool
```

## Usage

### Creating A Connection Pool

```go
package main

import (
	"fmt"
	"time"

	"github.com/meliadamian17/tcp-pool"
)

func main() {
	// Configure the connection pool
	config := pool.NewConfig(
		"localhost:9999",         // Address
		"test-pool",              // Pool name
		5,                        // Max connections
		2*time.Second,            // Connection timeout
		10*time.Second,           // Idle timeout
		3,                        // Max retries
		pool.NewExponentialBackoff(1, 10), // Backoff strategy
		pool.PoolHooks{},         // Optional hooks
	)

	// Create the connection pool
	p, err := pool.New(*config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create connection pool: %v", err))
	}

	// Acquire a connection
	conn, err := p.Get()
	if err != nil {
		panic(fmt.Sprintf("Failed to get connection: %v", err))
	}
	
	... Do something with connection

	// Release the connection back to the pool
	if err := p.Release(conn); err != nil {
		fmt.Printf("Failed to release connection: %v\n", err)
	}

	fmt.Println("Connection pool example completed successfully")
}

```
### Using Hooks
```go
hooks := pool.PoolHooks{
    OnConnectionCreate: func(conn net.Conn) {
        fmt.Println("Connection created:", conn.RemoteAddr())
    },
    OnConnectionAcquire: func(conn net.Conn) {
        fmt.Println("Connection acquired:", conn.RemoteAddr())
    },
    OnConnectionRelease: func(conn net.Conn) {
        fmt.Println("Connection released:", conn.RemoteAddr())
    },
    OnConnectionError: func(err error) {
        fmt.Println("Connection error:", err)
    },
}
config := pool.NewConfig(
    "localhost:9999",
    "hook-enabled-pool",
    5,
    2*time.Second,
    10*time.Second,
    3,
    pool.NewFibonacciBackoff(10),
    hooks,
)

```

## Contributing
Contributions are welcome! Please fork the repository, make your changes, and open a pull request.

## License
This project is licensed under the MIT License.

