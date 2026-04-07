package dto

type SetLineupRequest struct {
	PlayerIDs []string `json:"playerIds"`
}

type LineupResponse struct {
	GameID    string   `json:"gameId"`
	TeamID    string   `json:"teamId"`
	PlayerIDs []string `json:"playerIds"`
}
