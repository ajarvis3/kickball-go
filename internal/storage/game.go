package storage

import "github.com/ajarvis3/kickball-go/internal/domain"

type GameItem struct {
	PK           string           `dynamodbav:"PK"`
	SK           string           `dynamodbav:"SK"`
	GameID       string           `dynamodbav:"GameID"`
	LeagueID     string           `dynamodbav:"LeagueID"`
	RulesVersion int              `dynamodbav:"rulesVersion"`
	HomeTeamID   string           `dynamodbav:"homeTeamID"`
	AwayTeamID   string           `dynamodbav:"awayTeamID"`
	State        domain.GameState `dynamodbav:"state"`

	GSILeagueGamePK string `dynamodbav:"GSILeagueGamePK"`
	GSILeagueGameSK string `dynamodbav:"GSILeagueGameSK"`
}
