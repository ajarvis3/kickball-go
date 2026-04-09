package services

import (
	"github.com/ajarvis3/kickball-go/internal/data/domain"
	"github.com/ajarvis3/kickball-go/pkg/apperrors"
)

type RulesEngine struct {
}

func NewRulesEngine() *RulesEngine {
	return &RulesEngine{}
}

// ValidateAtBat enforces simple pitch/at-bat constraints from league rules.
// It ensures strikes, balls, and fouls are within configured maxima.
func (e *RulesEngine) ValidateAtBat(game domain.Game, rules domain.LeagueRules, atbat domain.AtBat) error {
	if rules.MaxStrikes > 0 && atbat.Strikes > rules.MaxStrikes {
		return apperrors.ErrInvalidInput
	}
	if rules.MaxBalls > 0 && atbat.Balls > rules.MaxBalls {
		return apperrors.ErrInvalidInput
	}
	if rules.MaxFouls > 0 && atbat.Fouls > rules.MaxFouls {
		return apperrors.ErrInvalidInput
	}

	return nil
}

// DoesInningMercyApply returns true when the mercy rule should be checked
// for the current inning according to the league rules.
func (e *RulesEngine) DoesInningMercyApply(rules domain.LeagueRules, game domain.Game, idx int) bool {
	applies := rules.MercyAppliesLastInning || (!rules.MercyAppliesLastInning && game.State.Inning == rules.MaxInnings)
	if !applies {
		return false
	}
	if idx < 0 || idx >= len(game.State.InningRuns) {
		return false
	}
	return game.State.InningRuns[idx] >= rules.MercyRunsPerInning
}
