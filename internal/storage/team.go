package storage

type TeamItem struct {
	PK       string `dynamodbav:"PK"`
	SK       string `dynamodbav:"SK"`
	TeamID   string `dynamodbav:"TeamID"`
	LeagueID string `dynamodbav:"LeagueID"`
	Name     string `dynamodbav:"name"`
}
