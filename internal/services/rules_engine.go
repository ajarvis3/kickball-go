package services

import "github.com/ajarvis3/kickball-go/internal/domain"

type RulesEngine struct {
	// TODO: config, etc.
}

func NewRulesEngine() *RulesEngine {
	return &RulesEngine{}
}

func (e *RulesEngine) ValidateAtBat(game domain.Game, rules domain.LeagueRules, atbat domain.AtBat) error {
	// TODO: enforce strikes, fouls, balls, mercy rules
	return nil
}