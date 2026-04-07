package storage

type LineupItem struct {
	PK        string   `dynamodbav:"PK"`
	SK        string   `dynamodbav:"SK"`
	GameID    string   `dynamodbav:"GameID"`
	TeamID    string   `dynamodbav:"TeamID"`
	PlayerIDs []string `dynamodbav:"playerIds"`
}
