package tencho

import (
	"github.com/ascenttree/tengo/common"
)

func ReadAuthenticationRequest(stream *common.IOStream) *AuthenticationRequest {
	username := stream.ReadString()

	return &AuthenticationRequest{
		Username: username,
	}
}

func ReadMatchStateChangeRequest(stream *common.IOStream) *MatchStateChangeRequest {
	newState := stream.ReadU8()

	return &MatchStateChangeRequest{
		NewState: newState,
	}
}

func ReadBeatmap(stream *common.IOStream) *Beatmap {
	filename := stream.ReadString()

	return &Beatmap{
		Filename: filename,
	}
}

func ReadTrackingPoint(stream *common.IOStream) *TrackingPoint {
	basePosX := stream.ReadF32()
	basePosY := stream.ReadF32()
	windowDeltaX := stream.ReadF32()
	windowDeltaY := stream.ReadF32()

	return &TrackingPoint{
		BasePosX:     basePosX,
		BasePosY:     basePosY,
		WindowDeltaX: windowDeltaX,
		WindowDeltaY: windowDeltaY,
	}
}

func ReadPlayerData(stream *common.IOStream) *PlayerData {
	username := stream.ReadString()
	state := stream.ReadU8()

	trackingPoints := []*TrackingPoint{}
	length := stream.ReadU32()
	for i := 0; i < int(length); i++ {
		trackingPoints = append(trackingPoints, ReadTrackingPoint(stream))
	}

	return &PlayerData{
		Username:       username,
		State:          state,
		TrackingPoints: trackingPoints,
	}
}
