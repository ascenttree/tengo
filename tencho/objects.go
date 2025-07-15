package tencho

type AuthenticationRequest struct {
	Username string
}

type MatchStateChangeRequest struct {
	NewState uint8
}

type Beatmap struct {
	Filename string
}

type TrackingPoint struct {
	BasePosX     float32
	BasePosY     float32
	WindowDeltaX float32
	WindowDeltaY float32
}

type PlayerData struct {
	Username       string
	State          uint8
	TrackingPoints []*TrackingPoint
}

type MatchData struct {
	MatchID    uint32
	MatchState uint8
	Beatmap    *Beatmap
	Players    map[string]*PlayerData
}
