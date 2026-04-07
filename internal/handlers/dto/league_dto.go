package dto

type CreateLeagueRequest struct {
	Name string `json:"name"`
}

type LeagueResponse struct {
	LeagueID            string `json:"leagueId"`
	Name                string `json:"name"`
	CurrentRulesVersion int    `json:"currentRulesVersion"`
}
