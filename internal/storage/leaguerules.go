package storage

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
)

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

func LeagueRulesToItem(r domain.LeagueRules) LeagueRulesItem {
	return LeagueRulesItem{
		PK:                     keys.LeaguePK(r.LeagueID),
		SK:                     keys.RulesSK(r.RulesVersion),
		LeagueID:               r.LeagueID,
		RulesVersion:           r.RulesVersion,
		MaxStrikes:             r.MaxStrikes,
		MaxBalls:               r.MaxBalls,
		MaxFouls:               r.MaxFouls,
		MaxInnings:             r.MaxInnings,
		MercyRunsPerInning:     r.MercyRunsPerInning,
		MercyAppliesLastInning: r.MercyAppliesLastInning,
		GameMercyRuns:          r.GameMercyRuns,
		Metadata:               r.Metadata,
	}
}

func ItemToLeagueRules(it LeagueRulesItem) domain.LeagueRules {
	return domain.LeagueRules{
		LeagueID:               it.LeagueID,
		RulesVersion:           it.RulesVersion,
		MaxStrikes:             it.MaxStrikes,
		MaxBalls:               it.MaxBalls,
		MaxFouls:               it.MaxFouls,
		MaxInnings:             it.MaxInnings,
		MercyRunsPerInning:     it.MercyRunsPerInning,
		MercyAppliesLastInning: it.MercyAppliesLastInning,
		GameMercyRuns:          it.GameMercyRuns,
		Metadata:               it.Metadata,
	}
}
