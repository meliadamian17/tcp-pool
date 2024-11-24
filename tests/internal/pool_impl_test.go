package internal

import (
	"fmt"
	"testing"
	"time"

	"github.com/meliadamian17/tcp-pool/internal"
	"github.com/meliadamian17/tcp-pool/tests/utils"
)

func TestPoolInitialization(t *testing.T) {

	serverConfig := utils.MockServerConfig{
		SendData: false,
	}
	server, address := utils.NewMockServer(t, serverConfig)
	defer server.Stop()

	config := internal.ConfigImpl{
		Address:        address,
		MaxConnections: 5,
		ConnTimeout:    2 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxRetries:     3,
		Backoff:        &utils.MockBackoff{},
	}

	pool, err := internal.NewConnectionPool(config)

	utils.AssertNil(t, err, "Pool initialization should not return an error")
	utils.AssertNotNil(t, pool, "Pool should be created successfully")
}

func TestPoolIdleTimeout(t *testing.T) {

	serverConfig := utils.MockServerConfig{
		SendData: false,
	}
	server, address := utils.NewMockServer(t, serverConfig)
	defer server.Stop()

	config := internal.ConfigImpl{
		Address:        address,
		MaxConnections: 5,
		ConnTimeout:    2 * time.Second,
		IdleTimeout:    1 * time.Second,
		MaxRetries:     3,
		Backoff:        &utils.MockBackoff{},
	}
	pool, _ := internal.NewConnectionPool(config)

	conn, _ := pool.Get()
	err := pool.Release(conn)
	utils.AssertNil(t, err, "Releasing a connection should not return an error")

	initialIdleConns := len(pool.IdleConns)
	fmt.Printf("IdleConns before cleanup: %d\n", initialIdleConns)
	utils.AssertEqual(t, 1, initialIdleConns, "There should be 1 idle connection before cleanup")

	time.Sleep(3 * time.Second)

	finalIdleConns := len(pool.IdleConns)
	fmt.Printf("IdleConns after cleanup: %d\n", finalIdleConns)
	utils.AssertEqual(t, 0, finalIdleConns, "Idle connection should have been cleaned up")
}
