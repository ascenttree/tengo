package tencho

import (
	"fmt"
	"sync"
)

// Player collection
type PlayerCollection struct {
	mu      sync.RWMutex
	Players map[string]*Player
}

func NewPlayerCollection() *PlayerCollection {
	return &PlayerCollection{
		Players: make(map[string]*Player),
	}
}

func (pc *PlayerCollection) Add(player *Player) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.Players[player.Username] = player
}

func (pc *PlayerCollection) Remove(player *Player) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	delete(pc.Players, player.Username)
	if player.Match != nil {
		player.Match.RemovePlayer(player)
	}
}

func (pc *PlayerCollection) AsList() []*Player {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	playersList := make([]*Player, 0, len(pc.Players))

	for _, player := range pc.Players {
		playersList = append(playersList, player)
	}

	return playersList
}

func (pc *PlayerCollection) SendPacket(packetId uint16, data []byte) error {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	for _, player := range pc.Players {
		err := player.SendPacket(packetId, data)
		if err != nil {
			return fmt.Errorf("failed to send packet: %v", err)
		}
	}
	return nil
}

// Match collection
type MatchCollection struct {
	mu      sync.RWMutex
	nextId  int
	Matches map[uint32]*Match
}

func NewMatchCollection() *MatchCollection {
	return &MatchCollection{
		Matches: make(map[uint32]*Match),
		nextId:  1,
	}
}

func (mc *MatchCollection) Add(match *Match) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.Matches[match.MatchID] = match
}

func (mc *MatchCollection) Remove(match *Match) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.Matches, match.MatchID)
	for _, player := range match.Players {
		player.Match = nil
	}

	if len(mc.Matches) == 0 {
		mc.nextId = 1
	}
}

func (mc *MatchCollection) AsList() []*Match {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	matchList := make([]*Match, 0, len(mc.Matches))

	for _, match := range mc.Matches {
		matchList = append(matchList, match)
	}

	return matchList
}

func (mc *MatchCollection) Matchmake(player *Player) (*Match, error) {
	// Normally, we would lock here, but doing so would create a deadlock with AddMatch (I found that out, the very fucking hard way)

	for _, match := range mc.Matches {
		if match.MatchState == MATCH_STATE_GATHERING_PLAYERS || match.MatchState == MATCH_STATE_RESULTS {
			return match, match.AddPlayer(player)
		}
	}

	matchId := mc.nextId
	mc.nextId++

	newMatch := NewMatch(matchId)
	err := newMatch.AddPlayer(player)
	if err != nil {
		return nil, err
	}

	mc.Add(newMatch)
	return newMatch, nil
}
