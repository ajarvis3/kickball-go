package domain

type Game struct {
	PK               string `dynamodbav:"PK"`
	SK               string `dynamodbav:"SK"`
	GameID           string `dynamodbav:"GameID"`
	LeagueID         string `dynamodbav:"LeagueID"`
	RulesVersion     int    `dynamodbav:"rulesVersion"`
	HomeTeamID       string `dynamodbav:"homeTeamID"`
	AwayTeamID       string `dynamodbav:"awayTeamID"`
	State            GameState

	GSILeagueGamePK string `dynamodbav:"GSILeagueGamePK"`
	GSILeagueGameSK string `dynamodbav:"GSILeagueGameSK"`
}

type GameState struct {
	Inning    int
	Half      string
	Outs      int
	HomeScore int
	AwayScore int
}