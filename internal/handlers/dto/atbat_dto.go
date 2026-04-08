package dto

type RecordAtBatRequest struct {
	GameID   string `json:"gameId"`
	LeagueID string `json:"leagueId"`
	PlayerID string `json:"playerId"`
	TeamID   string `json:"teamId"`

	Strikes int    `json:"strikes"`
	Balls   int    `json:"balls"`
	Fouls   int    `json:"fouls"`
	Result  string `json:"result"`
	RBI     int    `json:"rbi"`
}

type AtBatResponse struct {
	GameID   string `json:"gameId"`
	PlayerID string `json:"playerId"`
	TeamID   string `json:"teamId"`
	Seq      int    `json:"seq"`

	Inning  int    `json:"inning"`
	Half    string `json:"half"`
	Strikes int    `json:"strikes"`
	Balls   int    `json:"balls"`
	Fouls   int    `json:"fouls"`
	Result  string `json:"result"`
	RBI     int    `json:"rbi"`
}
