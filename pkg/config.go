package pool

import (
	"fmt"
	"io"
	"log"
	"meliadamian17/tcp-pool/utils"
	"os"
	"time"
)

type Config struct {
	address        string
	name           string
	maxConnections int
	connTimeout    time.Duration
	idleTimeout    time.Duration
	outputStream   *log.Logger
}

func NewConfig(
	address, name string,
	maxConnections int,
	connTimeout, idleTimeout time.Duration,
	outputStream io.Writer,
) *Config {
	if outputStream == nil {
		outputStream = os.Stdout
	}
	if len(name) == 0 {
		name = utils.GenerateNameByAddress(address)
	}
	return &Config{
		address:        address,
		name:           name,
		maxConnections: maxConnections,
		connTimeout:    connTimeout,
		idleTimeout:    idleTimeout,
		outputStream: log.New(
			outputStream,
			fmt.Sprintf("Connection Pool (%v, %v)", name, address),
			log.LstdFlags,
		),
	}
}
