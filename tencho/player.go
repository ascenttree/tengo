package tencho

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/ascenttree/tengo/common"
)

type Player struct {
	Conn         net.Conn
	Username     string
	LastPingTime int64
	Match        *Match
	MatchData    *PlayerData
	Server       *TenchoServer // Blame old me for not doing proper state management, fuck's sake
	Logger       *common.Logger
}

func (player *Player) SendPacket(packetId uint16, data []byte) error {
	stream := common.NewIOStream([]byte{}, binary.LittleEndian)

	stream.WriteU16(packetId)
	stream.WriteU8(0) // Compression flag (unused)
	stream.WriteU32(uint32(len(data)))
	stream.Push(data)

	_, err := player.Conn.Write(stream.Get())
	if err != nil {
		return fmt.Errorf("failed to send packet: %v", err)
	}

	return nil
}
