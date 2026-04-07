package storage

type AtbatItem struct {
	PK               string `dynamodbav:"PK"`
	SK               string `dynamodbav:"SK"`
	GameID           string `dynamodbav:"GameID"`
	LeagueID         string `dynamodbav:"LeagueID"`
	TeamID           string `dynamodbav:"TeamID"`
	PlayerID         string `dynamodbav:"PlayerID"`
	Seq              int    `dynamodbav:"seq"`
	Inning           int    `dynamodbav:"inning"`
	Half             string `dynamodbav:"half"`
	Strikes          int    `dynamodbav:"strikes"`
	Balls            int    `dynamodbav:"balls"`
	Fouls            int    `dynamodbav:"fouls"`
	Result           string `dynamodbav:"result"`
	RBI              int    `dynamodbav:"rbi"`
	Pitches          int    `dynamodbav:"pitches"`
	GSIPlayerAtBatPK string `dynamodbav:"GSIPlayerAtBatPK"`
	GSIPlayerAtBatSK string `dynamodbav:"GSIPlayerAtBatSK"`
}
