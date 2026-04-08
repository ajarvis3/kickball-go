package services

import (
	"errors"
	"testing"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/pkg/apperrors"
)

func makeGame(homeTeam, awayTeam string, inning int, half string, outs, homeScore, awayScore int, maxInnings int) domain.Game {
	return domain.Game{
		GameID:     "g1",
		LeagueID:   "l1",
		HomeTeamID: homeTeam,
		AwayTeamID: awayTeam,
		State: domain.GameState{
			Inning:     inning,
			Half:       half,
			Outs:       outs,
			HomeScore:  homeScore,
			AwayScore:  awayScore,
			InningRuns: make([]int, 2*maxInnings),
		},
	}
}

func makeRules(maxInnings int) domain.LeagueRules {
	return domain.LeagueRules{
		MaxInnings: maxInnings,
		MaxStrikes: 3,
		MaxBalls:   4,
	}
}

func TestApplyAtBat_NilEngine(t *testing.T) {
	var e *GameEngine
	_, err := e.ApplyAtBat(domain.Game{}, domain.LeagueRules{}, domain.AtBat{})
	if !errors.Is(err, apperrors.ErrInternalError) {
		t.Errorf("expected ErrInternalError, got %v", err)
	}
}

func TestApplyAtBat_NilRules(t *testing.T) {
	e := &GameEngine{Rules: nil}
	_, err := e.ApplyAtBat(domain.Game{}, domain.LeagueRules{}, domain.AtBat{})
	if !errors.Is(err, apperrors.ErrInternalError) {
		t.Errorf("expected ErrInternalError, got %v", err)
	}
}

func TestApplyAtBat_InvalidStrikeCount(t *testing.T) {
	e := NewGameEngine(NewRulesEngine())
	game := makeGame("ht", "at", 1, "top", 0, 0, 0, 7)
	rules := makeRules(7)
	atbat := domain.AtBat{TeamID: "ht", Strikes: 10, Result: "strikeout"}
	_, err := e.ApplyAtBat(game, rules, atbat)
	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestApplyAtBat_SingleWithRBI(t *testing.T) {
	e := NewGameEngine(NewRulesEngine())
	game := makeGame("ht", "at", 1, "top", 0, 0, 0, 7)
	rules := makeRules(7)
	atbat := domain.AtBat{TeamID: "ht", Result: "single", RBI: 2}
	updated, err := e.ApplyAtBat(game, rules, atbat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.State.HomeScore != 2 {
		t.Errorf("HomeScore = %d; want 2", updated.State.HomeScore)
	}
	if updated.State.InningRuns[0] != 2 {
		t.Errorf("InningRuns[0] = %d; want 2", updated.State.InningRuns[0])
	}
}

func TestApplyAtBat_OutIncrementsOuts(t *testing.T) {
	e := NewGameEngine(NewRulesEngine())
	game := makeGame("ht", "at", 1, "top", 0, 0, 0, 7)
	rules := makeRules(7)
	atbat := domain.AtBat{TeamID: "ht", Result: "out"}
	updated, err := e.ApplyAtBat(game, rules, atbat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.State.Outs != 1 {
		t.Errorf("Outs = %d; want 1", updated.State.Outs)
	}
}

func TestApplyAtBat_ThreeOutsAdvancesTopToBottom(t *testing.T) {
	e := NewGameEngine(NewRulesEngine())
	game := makeGame("ht", "at", 1, "top", 2, 0, 0, 7)
	rules := makeRules(7)
	atbat := domain.AtBat{TeamID: "ht", Result: "out"}
	updated, err := e.ApplyAtBat(game, rules, atbat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.State.Half != "bottom" {
		t.Errorf("Half = %q; want bottom", updated.State.Half)
	}
	if updated.State.Outs != 0 {
		t.Errorf("Outs = %d; want 0", updated.State.Outs)
	}
	if updated.State.Inning != 1 {
		t.Errorf("Inning = %d; want 1", updated.State.Inning)
	}
}

func TestApplyAtBat_ThreeOutsAtBottomAdvancesToNextInning(t *testing.T) {
	e := NewGameEngine(NewRulesEngine())
	game := makeGame("ht", "at", 1, "bottom", 2, 0, 0, 7)
	rules := makeRules(7)
	atbat := domain.AtBat{TeamID: "at", Result: "out"}
	updated, err := e.ApplyAtBat(game, rules, atbat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.State.Inning != 2 {
		t.Errorf("Inning = %d; want 2", updated.State.Inning)
	}
	if updated.State.Half != "top" {
		t.Errorf("Half = %q; want top", updated.State.Half)
	}
	if updated.State.Outs != 0 {
		t.Errorf("Outs = %d; want 0", updated.State.Outs)
	}
}

func TestApplyAtBat_DoublePlayAddsTwoOuts(t *testing.T) {
	e := NewGameEngine(NewRulesEngine())
	game := makeGame("ht", "at", 1, "top", 0, 0, 0, 7)
	rules := makeRules(7)
	atbat := domain.AtBat{TeamID: "ht", Result: "doubleplay"}
	updated, err := e.ApplyAtBat(game, rules, atbat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.State.Outs != 2 {
		t.Errorf("Outs = %d; want 2", updated.State.Outs)
	}
}

func TestApplyAtBat_TriplePlayAddsThreeOuts(t *testing.T) {
	e := NewGameEngine(NewRulesEngine())
	game := makeGame("ht", "at", 1, "top", 0, 0, 0, 7)
	rules := makeRules(7)
	atbat := domain.AtBat{TeamID: "ht", Result: "tripleplay"}
	updated, err := e.ApplyAtBat(game, rules, atbat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 3 outs triggers inning advance, so outs reset to 0
	if updated.State.Half != "bottom" {
		t.Errorf("Half = %q; want bottom after tripleplay", updated.State.Half)
	}
}

func TestApplyAtBat_GameMercyRuleTriggered(t *testing.T) {
	e := NewGameEngine(NewRulesEngine())
	rules := domain.LeagueRules{
		MaxInnings:    7,
		MaxStrikes:    3,
		MaxBalls:      4,
		GameMercyRuns: 10,
	}
	game := domain.Game{
		GameID:     "g1",
		LeagueID:   "l1",
		HomeTeamID: "ht",
		AwayTeamID: "at",
		State: domain.GameState{
			Inning:     3,
			Half:       "top",
			Outs:       0,
			HomeScore:  0,
			AwayScore:  0,
			InningRuns: make([]int, 14),
		},
	}
	// Home team scores 10 RBI
	atbat := domain.AtBat{TeamID: "ht", Result: "homerun", RBI: 10}
	updated, err := e.ApplyAtBat(game, rules, atbat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.State.Inning != rules.MaxInnings+1 {
		t.Errorf("Inning = %d; want %d (game over)", updated.State.Inning, rules.MaxInnings+1)
	}
	if updated.State.Half != "" {
		t.Errorf("Half = %q; want empty string", updated.State.Half)
	}
}

func TestApplyAtBat_AwayTeamScores(t *testing.T) {
	e := NewGameEngine(NewRulesEngine())
	game := makeGame("ht", "at", 1, "top", 0, 0, 0, 7)
	rules := makeRules(7)
	// idx for inning 1, top = 0
	atbat := domain.AtBat{TeamID: "at", Result: "double", RBI: 3}
	updated, err := e.ApplyAtBat(game, rules, atbat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.State.AwayScore != 3 {
		t.Errorf("AwayScore = %d; want 3", updated.State.AwayScore)
	}
}

func TestApplyAtBat_StrikeoutIncrementsOuts(t *testing.T) {
	e := NewGameEngine(NewRulesEngine())
	game := makeGame("ht", "at", 1, "top", 1, 0, 0, 7)
	rules := makeRules(7)
	atbat := domain.AtBat{TeamID: "ht", Result: "strikeout"}
	updated, err := e.ApplyAtBat(game, rules, atbat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.State.Outs != 2 {
		t.Errorf("Outs = %d; want 2", updated.State.Outs)
	}
}
