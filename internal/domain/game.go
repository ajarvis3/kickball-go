package domain

type Game struct {
	GameID       string
	LeagueID     string
	RulesVersion int
	HomeTeamID   string
	AwayTeamID   string
	State        GameState
}

type GameState struct {
	Inning    int
	Half      string
	Outs      int
	HomeScore int
	AwayScore int
}