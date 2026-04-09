package dto

type CreateTeamRequest struct {
    LeagueId string `json:"leagueId"`
    Name     string `json:"name"`
}

type TeamResponse struct {
    TeamID   string `json:"teamId"`
    LeagueID string `json:"leagueId"`
    Name     string `json:"name"`
}
