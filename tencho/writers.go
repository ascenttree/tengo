package tencho

import (
	"encoding/binary"

	"github.com/ascenttree/tengo/common"
)

func (beatmap *Beatmap) Serialize() []byte {
	stream := common.NewIOStream([]byte{}, binary.LittleEndian)
	stream.WriteString(beatmap.Filename)
	return stream.Get()
}

func (point *TrackingPoint) Serialize() []byte {
	stream := common.NewIOStream([]byte{}, binary.LittleEndian)

	stream.WriteF32(point.BasePosX)
	stream.WriteF32(point.BasePosY)
	stream.WriteF32(point.WindowDeltaX)
	stream.WriteF32(point.WindowDeltaY)

	return stream.Get()
}

func (player *PlayerData) Serialize() []byte {
	stream := common.NewIOStream([]byte{}, binary.LittleEndian)

	stream.WriteString(player.Username)
	stream.WriteU8(player.State)
	stream.WriteU32(uint32(len(player.TrackingPoints)))

	for _, trackingPoint := range player.TrackingPoints {
		stream.Write(trackingPoint.Serialize())
	}

	return stream.Get()
}

func (match *MatchData) Serialize() []byte {
	stream := common.NewIOStream([]byte{}, binary.LittleEndian)

	stream.WriteU32(match.MatchID)
	stream.WriteU8(match.MatchState)
	stream.Write(match.Beatmap.Serialize())
	stream.WriteU32(uint32(len(match.Players)))

	for _, player := range match.Players {
		stream.Write(player.Serialize())
	}

	return stream.Get()
}
