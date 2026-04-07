package domain

type LeagueRules struct {
	LeagueID     string `json:"leagueId"`
	RulesVersion int    `json:"rulesVersion"`

	// Pitching / at-bat rules
	MaxStrikes int `json:"maxStrikes"`
	MaxBalls   int `json:"maxBalls"`
	MaxFouls   int `json:"maxFouls"`

	// Inning rules
	MaxInnings             int  `json:"maxInnings"`
	MercyRunsPerInning     int  `json:"mercyRunsPerInning"`
	MercyAppliesLastInning bool `json:"mercyAppliesLastInning"`

	// Game-level mercy rule
	GameMercyRuns int `json:"gameMercyRuns"`

	// Future-proofing: allows adding new rule categories without breaking old versions
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
