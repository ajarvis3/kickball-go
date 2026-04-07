package mappers

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
)

func LeagueRulesToItem(r domain.LeagueRules) storage.LeagueRulesItem {
	return storage.LeagueRulesItem{
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

func ItemToLeagueRules(it storage.LeagueRulesItem) domain.LeagueRules {
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
