package services

import (
	"errors"
	"testing"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/pkg/apperrors"
)

func TestValidateAtBat_Valid(t *testing.T) {
	e := NewRulesEngine()
	game := domain.Game{}
	rules := domain.LeagueRules{MaxStrikes: 3, MaxBalls: 4, MaxFouls: 5}
	atbat := domain.AtBat{Strikes: 2, Balls: 3, Fouls: 4}
	if err := e.ValidateAtBat(game, rules, atbat); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateAtBat_TooManyStrikes(t *testing.T) {
	e := NewRulesEngine()
	rules := domain.LeagueRules{MaxStrikes: 3}
	atbat := domain.AtBat{Strikes: 4}
	err := e.ValidateAtBat(domain.Game{}, rules, atbat)
	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestValidateAtBat_TooManyBalls(t *testing.T) {
	e := NewRulesEngine()
	rules := domain.LeagueRules{MaxBalls: 4}
	atbat := domain.AtBat{Balls: 5}
	err := e.ValidateAtBat(domain.Game{}, rules, atbat)
	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestValidateAtBat_TooManyFouls(t *testing.T) {
	e := NewRulesEngine()
	rules := domain.LeagueRules{MaxFouls: 2}
	atbat := domain.AtBat{Fouls: 3}
	err := e.ValidateAtBat(domain.Game{}, rules, atbat)
	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestValidateAtBat_ZeroMaxMeansNoLimit(t *testing.T) {
	e := NewRulesEngine()
	rules := domain.LeagueRules{MaxStrikes: 0, MaxBalls: 0, MaxFouls: 0}
	atbat := domain.AtBat{Strikes: 100, Balls: 100, Fouls: 100}
	if err := e.ValidateAtBat(domain.Game{}, rules, atbat); err != nil {
		t.Errorf("expected no error when max is 0, got %v", err)
	}
}

func TestDoesInningMercyApply_True(t *testing.T) {
	e := NewRulesEngine()
	rules := domain.LeagueRules{
		MaxInnings:             7,
		MercyRunsPerInning:     5,
		MercyAppliesLastInning: true,
	}
	game := domain.Game{
		State: domain.GameState{
			Inning:     7,
			InningRuns: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 0},
		},
	}
	// idx for inning 7, top = (7-1)*2 = 12
	if !e.DoesInningMercyApply(rules, game, 12) {
		t.Errorf("expected mercy to apply")
	}
}

func TestDoesInningMercyApply_NotLastInning_WhenMercyDoesNotApplyToLastInning(t *testing.T) {
	e := NewRulesEngine()
	rules := domain.LeagueRules{
		MaxInnings:             7,
		MercyRunsPerInning:     5,
		MercyAppliesLastInning: false,
	}
	game := domain.Game{
		State: domain.GameState{
			Inning:     3,
			InningRuns: []int{0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	// applies only when inning == MaxInnings when MercyAppliesLastInning is false
	if e.DoesInningMercyApply(rules, game, 4) {
		t.Errorf("expected mercy not to apply for non-last inning")
	}
}

func TestDoesInningMercyApply_IdxOutOfBounds(t *testing.T) {
	e := NewRulesEngine()
	rules := domain.LeagueRules{MaxInnings: 7, MercyRunsPerInning: 5, MercyAppliesLastInning: true}
	game := domain.Game{State: domain.GameState{Inning: 7, InningRuns: []int{6}}}
	if e.DoesInningMercyApply(rules, game, 5) {
		t.Errorf("expected false for out-of-bounds idx")
	}
}

func TestDoesInningMercyApply_NegativeIdx(t *testing.T) {
	e := NewRulesEngine()
	rules := domain.LeagueRules{MaxInnings: 7, MercyRunsPerInning: 5, MercyAppliesLastInning: true}
	game := domain.Game{State: domain.GameState{Inning: 7, InningRuns: []int{6}}}
	if e.DoesInningMercyApply(rules, game, -1) {
		t.Errorf("expected false for negative idx")
	}
}

func TestDoesInningMercyApply_BelowThreshold(t *testing.T) {
	e := NewRulesEngine()
	rules := domain.LeagueRules{MaxInnings: 7, MercyRunsPerInning: 5, MercyAppliesLastInning: true}
	game := domain.Game{
		State: domain.GameState{
			Inning:     7,
			InningRuns: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0},
		},
	}
	if e.DoesInningMercyApply(rules, game, 12) {
		t.Errorf("expected false when runs below threshold")
	}
}
