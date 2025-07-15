package main

import (
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/ascenttree/tengo/common"
	"github.com/ascenttree/tengo/discovery"
	"github.com/ascenttree/tengo/tencho"
)

type Config struct {
	UDPServer struct {
		Port uint16
	}
	TCPServer struct {
		Port uint16
	}
}

func main() {
	config := &Config{}

	_, err := toml.DecodeFile("config.toml", config)
	if err != nil {
		fmt.Println("Error: could not decode config, please check config.toml and restart the server")
		return
	}

	discoveryServer := discovery.NewDiscoveryServer(
		config.UDPServer.Port,
	)

	tenchoServer := tencho.NewTenchoServer(
		"0.0.0.0",
		config.TCPServer.Port,
		common.CreateLogger("tencho", common.INFO),
	)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		discoveryServer.Serve()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tenchoServer.Serve()
	}()

	wg.Wait()
}
