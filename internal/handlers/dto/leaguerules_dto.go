package dto

type CreateLeagueRulesRequest struct {
	LeagueID               string                 `json:"leagueId"`
	RulesVersion           int                    `json:"rulesVersion"`
	MaxStrikes             int                    `json:"maxStrikes"`
	MaxBalls               int                    `json:"maxBalls"`
	MaxFouls               int                    `json:"maxFouls"`
	MaxInnings             int                    `json:"maxInnings"`
	MercyRunsPerInning     int                    `json:"mercyRunsPerInning"`
	MercyAppliesLastInning bool                   `json:"mercyAppliesLastInning"`
	GameMercyRuns          int                    `json:"gameMercyRuns"`
	Metadata               map[string]interface{} `json:"metadata,omitempty"`
}

type LeagueRulesResponse struct {
	LeagueID               string                 `json:"leagueId"`
	RulesVersion           int                    `json:"rulesVersion"`
	MaxStrikes             int                    `json:"maxStrikes"`
	MaxBalls               int                    `json:"maxBalls"`
	MaxFouls               int                    `json:"maxFouls"`
	MaxInnings             int                    `json:"maxInnings"`
	MercyRunsPerInning     int                    `json:"mercyRunsPerInning"`
	MercyAppliesLastInning bool                   `json:"mercyAppliesLastInning"`
	GameMercyRuns          int                    `json:"gameMercyRuns"`
	Metadata               map[string]interface{} `json:"metadata,omitempty"`
}
