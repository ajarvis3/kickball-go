package domain

import "github.com/ajarvis3/kickball-go/internal/types"

type Game struct {
    GameID       string
    LeagueID     string
    RulesVersion int
    HomeTeamID   string
    AwayTeamID   string
    State        GameState
}

type GameState struct {
    Inning     int
    Half       types.Half
    Outs       int
    HomeScore  int
    AwayScore  int
    InningRuns []int
}
