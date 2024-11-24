package internal

import (
	"net"
	"testing"
	"time"

	"meliadamian17/tcp-pool/internal"
	"meliadamian17/tcp-pool/tests/utils"
)

func TestPoolHooks(t *testing.T) {
	hookTriggered := false

	hooks := internal.PoolHooks{
		OnConnectionCreate: func(conn net.Conn) {
			hookTriggered = true
		},
	}

	serverConfig := utils.MockServerConfig{
		SendData:     true,
		Data:         []byte("test data"),
		SendInterval: 1 * time.Second,
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
		Hooks:          hooks,
	}
	pool, _ := internal.NewConnectionPool(config)

	conn, err := pool.Get()
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	defer pool.Release(conn)

	utils.AssertTrue(t, hookTriggered, "OnConnectionCreate hook should be triggered")
}
