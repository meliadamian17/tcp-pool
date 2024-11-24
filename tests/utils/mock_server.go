package utils

import (
	"net"
	"testing"
	"time"
)

type MockServerConfig struct {
	SendData     bool
	Data         []byte
	SendInterval time.Duration
}

type MockServer struct {
	listener net.Listener
	config   MockServerConfig
	stopChan chan struct{}
}

func NewMockServer(t *testing.T, config MockServerConfig) (*MockServer, string) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}

	server := &MockServer{
		listener: listener,
		config:   config,
		stopChan: make(chan struct{}),
	}

	go server.run()
	return server, listener.Addr().String()
}

func (s *MockServer) Stop() {
	close(s.stopChan)
	s.listener.Close()
}

func (s *MockServer) run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.stopChan:
				return
			default:
				continue
			}
		}

		go s.handleClient(conn)
	}
}

func (s *MockServer) handleClient(conn net.Conn) {
	defer conn.Close()

	if s.config.SendData {
		ticker := time.NewTicker(s.config.SendInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_, err := conn.Write(s.config.Data)
				if err != nil {
					return
				}
			case <-s.stopChan:
				return
			}
		}
	} else {

		<-s.stopChan
	}
}
