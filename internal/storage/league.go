package storage

type LeagueItem struct {
	PK                  string `dynamodbav:"PK"`
	SK                  string `dynamodbav:"SK"`
	LeagueID            string `dynamodbav:"leagueId"`
	Name                string `dynamodbav:"name"`
	CurrentRulesVersion int    `dynamodbav:"currentRulesVersion"`
}
