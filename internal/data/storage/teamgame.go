package storage

type TeamGameItem struct {
    PK         string `dynamodbav:"PK"`
    SK         string `dynamodbav:"SK"`
    GameID     string `dynamodbav:"GameID"`
    TeamID     string `dynamodbav:"TeamID"`
    OpponentID string `dynamodbav:"OpponentID"`
    LeagueID   string `dynamodbav:"LeagueID"`
    Date       string `dynamodbav:"date"`
}
