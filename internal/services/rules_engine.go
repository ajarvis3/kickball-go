package services

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
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
