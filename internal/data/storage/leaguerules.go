package storage

type LeagueRulesItem struct {
    PK                     string                 `dynamodbav:"PK"`
    SK                     string                 `dynamodbav:"SK"`
    LeagueID               string                 `dynamodbav:"leagueId"`
    RulesVersion           int                    `dynamodbav:"rulesVersion"`
    MaxStrikes             int                    `dynamodbav:"maxStrikes"`
    MaxBalls               int                    `dynamodbav:"maxBalls"`
    MaxFouls               int                    `dynamodbav:"maxFouls"`
    MaxInnings             int                    `dynamodbav:"maxInnings"`
    MercyRunsPerInning     int                    `dynamodbav:"mercyRunsPerInning"`
    MercyAppliesLastInning bool                   `dynamodbav:"mercyAppliesLastInning"`
    GameMercyRuns          int                    `dynamodbav:"gameMercyRuns"`
    Metadata               map[string]interface{} `dynamodbav:"metadata,omitempty"`
}
