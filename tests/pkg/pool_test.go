package tcppool

import (
	"testing"
	"time"

	pool "github.com/meliadamian17/tcppool"
	"github.com/meliadamian17/tcppool/tests/utils"
)

func TestPoolCreation(t *testing.T) {

	config := pool.NewConfig(
		"localhost:9999",
		"test-pool",
		5,
		2*time.Second,
		10*time.Second,
		3,
		pool.NewExponentialBackoff(1, 10),
		pool.PoolHooks{},
	)
	pool, err := pool.New(*config)

	utils.AssertNil(t, err, "Pool creation should not return an error")
	utils.AssertNotNil(t, pool, "Pool should be successfully created")
}

func TestPoolGetRelease(t *testing.T) {
	// Mock server configuration: sends "test data" every 1 second
	serverConfig := utils.MockServerConfig{
		SendData:     true,
		Data:         []byte("test data"),
		SendInterval: 1 * time.Second,
	}

	server, address := utils.NewMockServer(t, serverConfig)
	defer server.Stop()

	config := pool.NewConfig(
		address,
		"test-pool",
		5,
		2*time.Second,
		10*time.Second,
		3,
		pool.NewExponentialBackoff(1, 10),
		pool.PoolHooks{},
	)
	p, _ := pool.New(*config)

	conn, err := p.Get()
	utils.AssertNil(t, err, "Getting a connection should not return an error")
	utils.AssertNotNil(t, conn, "Connection should not be nil")

	err = p.Release(conn)
	utils.AssertNil(t, err, "Releasing a connection should not return an error")
}
