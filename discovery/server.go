package discovery

import (
	"fmt"
	"net"
	"time"
)

type DiscoveryServer struct {
	Port uint16
}

func (server *DiscoveryServer) Serve() error {
	conn, err := net.Dial("udp", fmt.Sprintf("255.255.255.255:%d", server.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to discovery server: %w", err)
	}

	defer conn.Close()

	for {
		conn.Write([]byte("Tenchooooooo"))
		time.Sleep(time.Second)
	}
}

func NewDiscoveryServer(port uint16) *DiscoveryServer {
	return &DiscoveryServer{
		Port: port,
	}
}
