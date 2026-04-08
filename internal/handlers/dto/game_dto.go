package dto

type CreateGameRequest struct {
	LeagueID   string `json:"leagueId"`
	HomeTeamID string `json:"homeTeamId"`
	AwayTeamID string `json:"awayTeamId"`
	Date       string `json:"date"` // ISO8601
}

type GameResponse struct {
	GameID       string       `json:"gameId"`
	LeagueID     string       `json:"leagueId"`
	HomeTeamID   string       `json:"homeTeamId"`
	AwayTeamID   string       `json:"awayTeamId"`
	RulesVersion int          `json:"rulesVersion"`
	State        GameStateDTO `json:"state"`
}

type GameStateDTO struct {
	Inning    int    `json:"inning"`
	Half      string `json:"half"`
	Outs      int    `json:"outs"`
	HomeScore int    `json:"homeScore"`
	AwayScore int    `json:"awayScore"`
}
