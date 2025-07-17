package tencho

import (
	"fmt"
	"time"

	"github.com/ascenttree/tengo/common"
)

var Handlers = map[uint16]func(*common.IOStream, *Player) error{}

func EnsureInMatch(handler func(*common.IOStream, *Player) error) func(*common.IOStream, *Player) error {
	return func(stream *common.IOStream, player *Player) error {
		if player.Match == nil {
			player.Logger.Error("User tried to send match packet while not in match?")
			return fmt.Errorf("not in match")
		}

		return handler(stream, player)
	}
}

func HandlePing(stream *common.IOStream, player *Player) error {
	player.LastPingTime = time.Now().Unix()
	return nil
}

func HandleAuthenticationRequest(stream *common.IOStream, player *Player) error {
	authRequest := ReadAuthenticationRequest(stream)

	player.Username = authRequest.Username
	player.Logger.Info("-> Authenticated")
	player.Server.Players.Add(player)

	// The authentication response packet is empty, hence []byte{}
	err := player.SendPacket(SERVER_AUTHENTICATION_RESPONSE, []byte{})
	if err != nil {
		return fmt.Errorf("failed to send packet: %v", err)
	}

	return nil
}

func HandleMatchmakeRequest(stream *common.IOStream, player *Player) error {
	match, err := player.Server.Matches.Matchmake(player)
	if err != nil {
		return fmt.Errorf("failed to matchmake: %v", err)
	}

	player.SendPacket(SERVER_MATCHMAKE_RESPONSE, match.Serialize())
	player.Logger.Info("-> Entered lobby")

	return nil
}

func HandleRequestStateChange(stream *common.IOStream, player *Player) error {
	matchStateChangeRequest := ReadMatchStateChangeRequest(stream)
	return player.Match.RequestStateChange(matchStateChangeRequest.NewState, player)
}

func HandleRequestSongChange(stream *common.IOStream, player *Player) error {
	newBeatmap := ReadBeatmap(stream)
	player.Logger.Info("-> Chose beatmap " + newBeatmap.Filename)
	return player.Match.ChooseSong(newBeatmap)
}

func HandleInputUpdate(stream *common.IOStream, player *Player) error {
	playerData := ReadPlayerData(stream)
	return player.Match.UpdateInput(player, playerData)
}

func HandleLeave(stream *common.IOStream, player *Player) error {
	player.Logger.Info("-> Left match")
	return player.Match.RemovePlayer(player)
}

func init() {
	Handlers[CLIENT_PING] = HandlePing
	Handlers[CLIENT_AUTHENTICATION_REQUEST] = HandleAuthenticationRequest
	Handlers[CLIENT_MATCHMAKE_REQUEST] = HandleMatchmakeRequest
	Handlers[CLIENT_REQUEST_STATE_CHANGE] = EnsureInMatch(HandleRequestStateChange)
	Handlers[CLIENT_REQUEST_SONG_CHANGE] = EnsureInMatch(HandleRequestSongChange)
	Handlers[CLIENT_INPUT_UPDATE] = EnsureInMatch(HandleInputUpdate)
	Handlers[CLIENT_LEAVE_MATCH] = EnsureInMatch(HandleLeave)
}
