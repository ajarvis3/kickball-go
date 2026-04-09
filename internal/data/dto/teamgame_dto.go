package dto

type TeamGameResponse struct {
    GameID     string `json:"gameId"`
    TeamID     string `json:"teamId"`
    OpponentID string `json:"opponentId"`
    LeagueID   string `json:"leagueId"`
    Date       string `json:"date"`
}
