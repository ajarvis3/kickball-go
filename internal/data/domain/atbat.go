package domain

import "github.com/ajarvis3/kickball-go/internal/types"

type AtBat struct {
    GameID   string
    LeagueID string
    TeamID   string
    PlayerID string

    Seq     int
    Inning  int
    Half    types.Half
    Strikes int
    Balls   int
    Fouls   int
    Result  string
    RBI     int
    Pitches int
}
