package main

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

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

	errGroup, ctx := errgroup.WithContext(context.Background())

	errGroup.Go(func() error {
		return discoveryServer.Serve(ctx)
	})

	errGroup.Go(func() error {
		return tenchoServer.Serve(ctx)
	})

	if err := errGroup.Wait(); err != nil {
		fmt.Printf("Error occurred while running Tengo: %v\n", err)
	}
}
