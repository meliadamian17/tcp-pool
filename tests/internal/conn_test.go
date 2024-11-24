package internal

import (
	"net"
	"testing"
	"time"

	"meliadamian17/tcp-pool/internal"
	"meliadamian17/tcp-pool/tests/utils"
)

func TestValidateConnection_ValidConnection(t *testing.T) {

	serverConfig := utils.MockServerConfig{
		SendData:     true,
		Data:         []byte("test data"),
		SendInterval: 250 * time.Millisecond,
	}

	server, address := utils.NewMockServer(t, serverConfig)
	defer server.Stop()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("Failed to connect to mock server: %v", err)
	}
	defer conn.Close()

	utils.AssertTrue(t, internal.Validate(conn), "Connection should be valid")

	conn.Close()
	utils.AssertFalse(t, internal.Validate(conn), "Connection should not be valid after close")
}

func TestValidateConnection_InvalidConnection(t *testing.T) {

	serverConfig := utils.MockServerConfig{
		SendData: false,
	}

	server, address := utils.NewMockServer(t, serverConfig)
	defer server.Stop()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("Failed to connect to mock server: %v", err)
	}
	defer conn.Close()

	utils.AssertFalse(t, internal.Validate(conn), "Conn Should Be Invalid")
}
