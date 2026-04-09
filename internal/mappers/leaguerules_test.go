package mappers

import (
	"testing"

	"github.com/ajarvis3/kickball-go/internal/data/domain"
	"github.com/ajarvis3/kickball-go/internal/data/storage"
	"github.com/ajarvis3/kickball-go/internal/keys"
)

func TestLeagueRulesToItem(t *testing.T) {
	r := domain.LeagueRules{
		LeagueID:               "l1",
		RulesVersion:           1,
		MaxStrikes:             3,
		MaxBalls:               4,
		MaxFouls:               0,
		MaxInnings:             7,
		MercyRunsPerInning:     5,
		MercyAppliesLastInning: true,
		GameMercyRuns:          10,
		Metadata:               map[string]interface{}{"key": "value"},
	}
	it := LeagueRulesToItem(r)
	if it.PK != keys.LeaguePK("l1") {
		t.Errorf("PK = %q; want %q", it.PK, keys.LeaguePK("l1"))
	}
	if it.SK != keys.RulesSK(1) {
		t.Errorf("SK = %q; want %q", it.SK, keys.RulesSK(1))
	}
	if it.LeagueID != "l1" {
		t.Errorf("LeagueID = %q; want l1", it.LeagueID)
	}
	if it.MaxStrikes != 3 {
		t.Errorf("MaxStrikes = %d; want 3", it.MaxStrikes)
	}
	if it.MaxInnings != 7 {
		t.Errorf("MaxInnings = %d; want 7", it.MaxInnings)
	}
	if !it.MercyAppliesLastInning {
		t.Errorf("MercyAppliesLastInning should be true")
	}
	if it.GameMercyRuns != 10 {
		t.Errorf("GameMercyRuns = %d; want 10", it.GameMercyRuns)
	}
}

func TestItemToLeagueRules(t *testing.T) {
	it := storage.LeagueRulesItem{
		LeagueID:               "l2",
		RulesVersion:           2,
		MaxStrikes:             4,
		MaxBalls:               3,
		MaxFouls:               5,
		MaxInnings:             9,
		MercyRunsPerInning:     8,
		MercyAppliesLastInning: false,
		GameMercyRuns:          15,
	}
	r := ItemToLeagueRules(it)
	if r.LeagueID != "l2" {
		t.Errorf("LeagueID = %q; want l2", r.LeagueID)
	}
	if r.MaxStrikes != 4 {
		t.Errorf("MaxStrikes = %d; want 4", r.MaxStrikes)
	}
	if r.MaxInnings != 9 {
		t.Errorf("MaxInnings = %d; want 9", r.MaxInnings)
	}
	if r.MercyAppliesLastInning {
		t.Errorf("MercyAppliesLastInning should be false")
	}
}

func TestLeagueRulesRoundTrip(t *testing.T) {
	r := domain.LeagueRules{
		LeagueID:     "lrt",
		RulesVersion: 1,
		MaxStrikes:   3,
		MaxInnings:   7,
	}
	back := ItemToLeagueRules(LeagueRulesToItem(r))
	if back.LeagueID != r.LeagueID || back.MaxStrikes != r.MaxStrikes || back.MaxInnings != r.MaxInnings {
		t.Errorf("round-trip mismatch: %+v", back)
	}
}
