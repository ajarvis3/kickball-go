package mappers

import (
	"testing"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/handlers/dto"
)

func TestRecordAtBatRequestToDomain(t *testing.T) {
	r := dto.RecordAtBatRequest{
		PlayerID: "p1",
		TeamID:   "t1",
		Strikes:  2,
		Balls:    3,
		Fouls:    1,
		Result:   "single",
		RBI:      1,
		Pitches:  7,
	}
	a := RecordAtBatRequestToDomain(r, "g1", "l1")
	if a.GameID != "g1" {
		t.Errorf("GameID = %q; want g1", a.GameID)
	}
	if a.LeagueID != "l1" {
		t.Errorf("LeagueID = %q; want l1", a.LeagueID)
	}
	if a.PlayerID != "p1" {
		t.Errorf("PlayerID = %q; want p1", a.PlayerID)
	}
	if a.Strikes != 2 {
		t.Errorf("Strikes = %d; want 2", a.Strikes)
	}
	if a.Result != "single" {
		t.Errorf("Result = %q; want single", a.Result)
	}
}

func TestAtBatToResponse(t *testing.T) {
	a := domain.AtBat{
		GameID:   "g1",
		PlayerID: "p1",
		TeamID:   "t1",
		Seq:      3,
		Inning:   2,
		Half:     "bottom",
		Strikes:  1,
		Balls:    2,
		Fouls:    0,
		Result:   "out",
		RBI:      0,
		Pitches:  4,
	}
	resp := AtBatToResponse(a)
	if resp.GameID != "g1" {
		t.Errorf("GameID = %q; want g1", resp.GameID)
	}
	if resp.Seq != 3 {
		t.Errorf("Seq = %d; want 3", resp.Seq)
	}
	if resp.Result != "out" {
		t.Errorf("Result = %q; want out", resp.Result)
	}
}

func TestCreateTeamRequestToDomain(t *testing.T) {
	r := dto.CreateTeamRequest{Name: "Tigers"}
	tm := CreateTeamRequestToDomain(r, "t1", "l1")
	if tm.Name != "Tigers" {
		t.Errorf("Name = %q; want Tigers", tm.Name)
	}
	if tm.TeamID != "t1" {
		t.Errorf("TeamID = %q; want t1", tm.TeamID)
	}
	if tm.LeagueID != "l1" {
		t.Errorf("LeagueID = %q; want l1", tm.LeagueID)
	}
}

func TestTeamToResponse(t *testing.T) {
	tm := domain.Team{TeamID: "t1", LeagueID: "l1", Name: "Tigers"}
	resp := TeamToResponse(tm)
	if resp.TeamID != "t1" || resp.Name != "Tigers" {
		t.Errorf("TeamToResponse mismatch: %+v", resp)
	}
}

func TestCreatePlayerRequestToDomain(t *testing.T) {
	r := dto.CreatePlayerRequest{Name: "Alice", Number: 5, Position: "ss"}
	p := CreatePlayerRequestToDomain(r, "p1", "t1", "l1")
	if p.Name != "Alice" {
		t.Errorf("Name = %q; want Alice", p.Name)
	}
	if p.Number != 5 {
		t.Errorf("Number = %d; want 5", p.Number)
	}
	if p.PlayerID != "p1" {
		t.Errorf("PlayerID = %q; want p1", p.PlayerID)
	}
}

func TestPlayerToResponse(t *testing.T) {
	p := domain.Player{PlayerID: "p1", TeamID: "t1", LeagueID: "l1", Name: "Alice", Number: 5, Position: "ss"}
	resp := PlayerToResponse(p)
	if resp.PlayerID != "p1" || resp.Name != "Alice" {
		t.Errorf("PlayerToResponse mismatch: %+v", resp)
	}
}

func TestCreateLeagueRequestToDomain(t *testing.T) {
	r := dto.CreateLeagueRequest{Name: "My League"}
	l := CreateLeagueRequestToDomain(r, "l1")
	if l.Name != "My League" {
		t.Errorf("Name = %q; want My League", l.Name)
	}
	if l.LeagueID != "l1" {
		t.Errorf("LeagueID = %q; want l1", l.LeagueID)
	}
	if l.CurrentRulesVersion != 1 {
		t.Errorf("CurrentRulesVersion = %d; want 1", l.CurrentRulesVersion)
	}
}

func TestLeagueToResponse(t *testing.T) {
	l := domain.League{LeagueID: "l1", Name: "My League", CurrentRulesVersion: 1}
	resp := LeagueToResponse(l)
	if resp.LeagueID != "l1" || resp.Name != "My League" {
		t.Errorf("LeagueToResponse mismatch: %+v", resp)
	}
}

func TestCreateLeagueRulesRequestToDomain(t *testing.T) {
	r := dto.CreateLeagueRulesRequest{
		MaxStrikes: 3,
		MaxBalls:   4,
		MaxInnings: 7,
	}
	rules := CreateLeagueRulesRequestToDomain(r, "l1", 2)
	if rules.LeagueID != "l1" {
		t.Errorf("LeagueID = %q; want l1", rules.LeagueID)
	}
	if rules.RulesVersion != 2 {
		t.Errorf("RulesVersion = %d; want 2", rules.RulesVersion)
	}
	if rules.MaxStrikes != 3 {
		t.Errorf("MaxStrikes = %d; want 3", rules.MaxStrikes)
	}
}

func TestLeagueRulesToResponse(t *testing.T) {
	r := domain.LeagueRules{LeagueID: "l1", RulesVersion: 1, MaxStrikes: 3, MaxInnings: 7}
	resp := LeagueRulesToResponse(r)
	if resp.LeagueID != "l1" || resp.MaxStrikes != 3 {
		t.Errorf("LeagueRulesToResponse mismatch: %+v", resp)
	}
}

func TestSetLineupRequestToDomain(t *testing.T) {
	r := dto.SetLineupRequest{PlayerIDs: []string{"p1", "p2"}}
	l := SetLineupRequestToDomain(r, "g1", "t1")
	if l.GameID != "g1" {
		t.Errorf("GameID = %q; want g1", l.GameID)
	}
	if len(l.PlayerIDs) != 2 {
		t.Errorf("PlayerIDs len = %d; want 2", len(l.PlayerIDs))
	}
}

func TestLineupToResponse(t *testing.T) {
	l := domain.Lineup{GameID: "g1", TeamID: "t1", PlayerIDs: []string{"p1"}}
	resp := LineupToResponse(l)
	if resp.GameID != "g1" || len(resp.PlayerIDs) != 1 {
		t.Errorf("LineupToResponse mismatch: %+v", resp)
	}
}

func TestGameToResponse(t *testing.T) {
	g := domain.Game{
		GameID:       "g1",
		LeagueID:     "l1",
		HomeTeamID:   "ht1",
		AwayTeamID:   "at1",
		RulesVersion: 1,
		State: domain.GameState{
			Inning:    2,
			Half:      "top",
			Outs:      1,
			HomeScore: 3,
			AwayScore: 1,
		},
	}
	resp := GameToResponse(g)
	if resp.GameID != "g1" {
		t.Errorf("GameID = %q; want g1", resp.GameID)
	}
	if resp.State.Inning != 2 {
		t.Errorf("State.Inning = %d; want 2", resp.State.Inning)
	}
	if resp.State.HomeScore != 3 {
		t.Errorf("State.HomeScore = %d; want 3", resp.State.HomeScore)
	}
}

func TestCreateGameRequestToDomain(t *testing.T) {
	r := dto.CreateGameRequest{HomeTeamID: "ht1", AwayTeamID: "at1"}
	g := CreateGameRequestToDomain(r, "g1", "l1")
	if g.GameID != "g1" {
		t.Errorf("GameID = %q; want g1", g.GameID)
	}
	if g.HomeTeamID != "ht1" {
		t.Errorf("HomeTeamID = %q; want ht1", g.HomeTeamID)
	}
	if g.RulesVersion != 1 {
		t.Errorf("RulesVersion = %d; want 1", g.RulesVersion)
	}
}

func TestTeamGameToResponse(t *testing.T) {
	tg := domain.TeamGame{GameID: "g1", TeamID: "t1", OpponentID: "t2", LeagueID: "l1", Date: "2024-01-01"}
	resp := TeamGameToResponse(tg)
	if resp.GameID != "g1" || resp.OpponentID != "t2" {
		t.Errorf("TeamGameToResponse mismatch: %+v", resp)
	}
}
