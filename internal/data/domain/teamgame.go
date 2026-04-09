package domain

type TeamGame struct {
    GameID     string
    TeamID     string
    OpponentID string
    LeagueID   string
    Date       string // ISO8601
}
