package dto

type RecordAtBatRequest struct {
    PlayerID string `json:"playerId"`
    TeamID   string `json:"teamId"`

    Strikes int    `json:"strikes"`
    Balls   int    `json:"balls"`
    Fouls   int    `json:"fouls"`
    Result  string `json:"result"`
    RBI     int    `json:"rbi"`
    Pitches int    `json:"pitches"`
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
    Pitches int    `json:"pitches"`
}