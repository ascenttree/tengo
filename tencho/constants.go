package tencho

const TENCHO_HEADER_SIZE int = 7

const (
	CLIENT_PING                    uint16 = 4
	SERVER_PING                    uint16 = 8
	SERVER_MATCH_DATA_UPDATE       uint16 = 26
	CLIENT_AUTHENTICATION_REQUEST  uint16 = 92
	SERVER_AUTHENTICATION_RESPONSE uint16 = 93
	CLIENT_MATCHMAKE_REQUEST       uint16 = 94
	SERVER_MATCHMAKE_RESPONSE      uint16 = 95
	CLIENT_REQUEST_STATE_CHANGE    uint16 = 96
	CLIENT_REQUEST_SONG_CHANGE     uint16 = 97
	SERVER_MATCH_PLAYER_UPDATE     uint16 = 98
	CLIENT_INPUT_UPDATE            uint16 = 99
	CLIENT_LEAVE_MATCH             uint16 = 100
)

const (
	MATCH_STATE_GATHERING_PLAYERS uint8 = iota
	MATCH_STATE_SONG_SELECT       uint8 = iota
	MATCH_STATE_DIFFICULTY_SELECT uint8 = iota
	MATCH_STATE_PREPARING         uint8 = iota
	MATCH_STATE_PLAYING           uint8 = iota
	MATCH_STATE_RESULTS           uint8 = iota
)

const (
	PLAYER_STATE_WAITING             uint8 = iota
	PLAYER_STATE_DIFFICULTY_SELECTED uint8 = iota
	PLAYER_STATE_READY_TO_PLAY       uint8 = iota
	PLAYER_STATE_PLAYING             uint8 = iota
	PLAYER_STATE_FINISHED            uint8 = iota
)
