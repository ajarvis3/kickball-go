package domain

type Player struct {
	PK       string `dynamodbav:"PK"`
	SK       string `dynamodbav:"SK"`
	PlayerID string `dynamodbav:"PlayerID"`
	TeamID   string `dynamodbav:"TeamID"`
	LeagueID string `dynamodbav:"LeagueID"`
	Name     string `dynamodbav:"name"`
	Number   int    `dynamodbav:"number"`
	Position string `dynamodbav:"position"`
}