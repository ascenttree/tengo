package tencho

import (
	"fmt"
	"sync"

	"github.com/ascenttree/tengo/common"
)

type Match struct {
	*MatchData
	mu      sync.RWMutex
	Players map[string]*Player
	Logger  *common.Logger
}

func NewMatch(matchId int) *Match {
	return &Match{
		MatchData: &MatchData{
			MatchID:    uint32(matchId),
			MatchState: MATCH_STATE_GATHERING_PLAYERS,
			Beatmap: &Beatmap{
				Filename: "",
			},
			Players: make(map[string]*PlayerData, 0),
		},
		Players: make(map[string]*Player, 0),
		Logger: common.CreateLogger(
			fmt.Sprintf("match %d", matchId),
			common.INFO,
		),
	}
}

func (match *Match) SendUpdates(packetId uint16) error {
	data := match.Serialize()

	for _, p := range match.Players {
		// seriously peppy? why are there 2 update packets
		err := p.SendPacket(packetId, data)
		if err != nil {
			return fmt.Errorf("failed to send match update: %v", err)
		}
	}

	return nil
}

func (match *Match) AddPlayer(player *Player) error {
	match.mu.Lock()
	defer match.mu.Unlock()

	player.Match = match
	playerData := &PlayerData{
		Username:       player.Username,
		State:          PLAYER_STATE_WAITING,
		TrackingPoints: []*TrackingPoint{},
	}
	match.Players[player.Username] = player
	match.MatchData.Players[player.Username] = playerData

	return match.SendUpdates(SERVER_MATCH_PLAYER_UPDATE)
}

func (match *Match) RemovePlayer(player *Player) error {
	match.mu.Lock()
	defer match.mu.Unlock()

	player.Match = nil
	delete(match.Players, player.Username)
	delete(match.MatchData.Players, player.Username)

	if len(match.Players) == 0 {
		match.Logger.Info("Match is empty, so it is being disposed")
		player.Server.Matches.Remove(match)
		return nil
	}

	return match.SendUpdates(SERVER_MATCH_PLAYER_UPDATE)
}

func (match *Match) RequestStateChange(newState uint8, player *Player) error {
	match.mu.Lock()
	defer match.mu.Unlock()

	switch newState {
	case MATCH_STATE_SONG_SELECT:
		if match.MatchState != MATCH_STATE_GATHERING_PLAYERS && match.MatchState != MATCH_STATE_RESULTS {
			return nil
		}

		match.Logger.Info("-> Selecting song")
		match.MatchState = MATCH_STATE_SONG_SELECT
	case MATCH_STATE_PREPARING:
		player.Logger.Info("-> Selected their difficulty")
		match.MatchData.Players[player.Username].State = PLAYER_STATE_DIFFICULTY_SELECTED
		match.CheckReady()
	case MATCH_STATE_PLAYING:
		match.MatchData.Players[player.Username].State = PLAYER_STATE_READY_TO_PLAY
		match.CheckReady()
	case MATCH_STATE_RESULTS:
		match.MatchData.Players[player.Username].State = PLAYER_STATE_FINISHED
		match.CheckReady()
	}

	return match.SendUpdates(SERVER_MATCH_DATA_UPDATE)
}

func (match *Match) ChooseSong(beatmap *Beatmap) error {
	if match.MatchState != MATCH_STATE_SONG_SELECT && match.MatchState != MATCH_STATE_DIFFICULTY_SELECT {
		return nil
	}

	match.Beatmap = beatmap
	if match.Beatmap.Filename == "" {
		match.MatchState = MATCH_STATE_SONG_SELECT
	} else {
		match.MatchState = MATCH_STATE_DIFFICULTY_SELECT
	}

	return match.SendUpdates(SERVER_MATCH_DATA_UPDATE)
}

func (match *Match) UpdateInput(player *Player, playerData *PlayerData) error {
	match.mu.Lock()
	defer match.mu.Unlock()

	match.MatchData.Players[player.Username].TrackingPoints = playerData.TrackingPoints
	return match.SendUpdates(SERVER_MATCH_PLAYER_UPDATE)
}

func (match *Match) CheckReady() {
	// We would lock, but that creates a deadlock with RequestStateChange

	switch match.MatchState {
	case MATCH_STATE_DIFFICULTY_SELECT:
		for _, player := range match.MatchData.Players {
			if player.State != PLAYER_STATE_DIFFICULTY_SELECTED {
				return
			}
		}

		match.MatchState = MATCH_STATE_PREPARING
	case MATCH_STATE_PREPARING:
		for _, player := range match.MatchData.Players {
			if player.State != PLAYER_STATE_READY_TO_PLAY {
				return
			}
		}

		match.Logger.Info("-> Game started")
		match.MatchState = MATCH_STATE_PLAYING
	case MATCH_STATE_PLAYING:
		for _, player := range match.MatchData.Players {
			if player.State != PLAYER_STATE_FINISHED {
				return
			}
		}

		match.Logger.Info("-> Game finished")
		match.MatchState = MATCH_STATE_RESULTS
	}

	match.SendUpdates(SERVER_MATCH_DATA_UPDATE)
}
