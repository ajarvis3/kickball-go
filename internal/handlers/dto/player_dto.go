package dto

type CreatePlayerRequest struct {
    Name     string `json:"name"`
    Number   int    `json:"number"`
    Position string `json:"position"`
}

type PlayerResponse struct {
    PlayerID string `json:"playerId"`
    TeamID   string `json:"teamId"`
    LeagueID string `json:"leagueId"`
    Name     string `json:"name"`
    Number   int    `json:"number"`
    Position string `json:"position"`
}