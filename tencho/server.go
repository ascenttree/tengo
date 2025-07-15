package tencho

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/ascenttree/tengo/common"
)

type TenchoServer struct {
	Host    string
	Port    uint16
	Logger  *common.Logger
	Matches *MatchCollection
	Players *PlayerCollection
}

func (server *TenchoServer) Serve() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port))
	if err != nil {
		server.Logger.Error(fmt.Sprintf("Failed to start Tencho server on %s:%d: %v", server.Host, server.Port, err))
		return
	}

	server.Logger.Info(fmt.Sprintf("Tencho server started on %s:%d", server.Host, server.Port))

	server.Players = NewPlayerCollection()
	server.Matches = NewMatchCollection()
	go server.StartPing()

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			server.Logger.Error(fmt.Sprintf("Failed to accept connection: %v", err))
			continue
		}

		go server.HandleConnection(conn)
	}
}

func (server *TenchoServer) StartPing() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, player := range server.Players.AsList() {
				if player.Conn != nil {
					// Send outgoing pings
					err := player.SendPacket(SERVER_PING, []byte{})
					if err != nil {
						player.Logger.Error(fmt.Sprintf("Failed to send ping to player %s: %v", player.Username, err))
					}

					// Check incoming pings
					if time.Now().Unix()-player.LastPingTime > 10 {
						player.Logger.Warning(fmt.Sprintf("Player %s has not responded to ping, disconnecting", player.Username))
						server.Players.RemovePlayer(player)
					}
				}
			}
		}
	}
}

func (server *TenchoServer) HandleConnection(conn net.Conn) {
	defer conn.Close()

	server.Logger.Info(fmt.Sprintf("Accepted connection from %s", conn.RemoteAddr()))

	logger := common.CreateLogger(
		conn.RemoteAddr().String(),
		server.Logger.GetLevel(),
	)

	player := &Player{
		Conn:         conn,
		Logger:       logger,
		LastPingTime: time.Now().Unix(),
		Server:       server,
	}

	player.Logger.Debug("-> Connected")

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			player.Logger.Info("-> Disconnected")
			server.Players.RemovePlayer(player)
			return
		}

		if n < TENCHO_HEADER_SIZE {
			server.Logger.Warning("Received data smaller than header size, ignoring")
			return
		}

		buffer = buffer[:n]
		headerReader := common.NewIOStream(buffer, binary.LittleEndian)

		for headerReader.Available() >= TENCHO_HEADER_SIZE {
			packetId := headerReader.ReadU16()
			headerReader.Skip(1) // Skip compression flag (unused in modern clients)
			packetSize := headerReader.ReadU32()

			packetData := headerReader.Read(int(packetSize))
			reader := common.NewIOStream(packetData, binary.LittleEndian)

			handler, ok := Handlers[packetId]
			if !ok {
				player.Logger.Error(fmt.Sprintf("Unknown packet ID: %d", packetId))
				continue
			}

			err = handler(reader, player)
			if err != nil {
				player.Logger.Error(fmt.Sprintf("Error handling packet ID %d: %v", packetId, err))
				continue
			}
		}
	}
}

func NewTenchoServer(host string, port uint16, logger *common.Logger) *TenchoServer {
	return &TenchoServer{
		Host:   host,
		Port:   port,
		Logger: logger,
	}
}
