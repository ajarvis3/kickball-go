package domain

type League struct {
	PK                  string `dynamodbav:"PK"`
	SK                  string `dynamodbav:"SK"`
	LeagueID            string `dynamodbav:"LeagueID"`
	Name                string `dynamodbav:"name"`
	CurrentRulesVersion int    `dynamodbav:"currentRulesVersion"`
}