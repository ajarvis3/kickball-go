package mappers

import (
	"testing"

	"github.com/ajarvis3/kickball-go/internal/data/domain"
	"github.com/ajarvis3/kickball-go/internal/data/storage"
	"github.com/ajarvis3/kickball-go/internal/keys"
)

func TestLeagueToItem(t *testing.T) {
	l := domain.League{
		LeagueID:            "l1",
		Name:                "Test League",
		CurrentRulesVersion: 2,
	}
	it := LeagueToItem(l)
	if it.PK != keys.LeaguePK("l1") {
		t.Errorf("PK = %q; want %q", it.PK, keys.LeaguePK("l1"))
	}
	if it.SK != keys.LeagueSK("l1") {
		t.Errorf("SK = %q; want %q", it.SK, keys.LeagueSK("l1"))
	}
	if it.LeagueID != "l1" {
		t.Errorf("LeagueID = %q; want l1", it.LeagueID)
	}
	if it.Name != "Test League" {
		t.Errorf("Name = %q; want Test League", it.Name)
	}
	if it.CurrentRulesVersion != 2 {
		t.Errorf("CurrentRulesVersion = %d; want 2", it.CurrentRulesVersion)
	}
}

func TestItemToLeague(t *testing.T) {
	it := storage.LeagueItem{
		PK:                  "LEAGUE#l2",
		SK:                  "LEAGUE#l2",
		LeagueID:            "l2",
		Name:                "Other League",
		CurrentRulesVersion: 3,
	}
	l := ItemToLeague(it)
	if l.LeagueID != "l2" {
		t.Errorf("LeagueID = %q; want l2", l.LeagueID)
	}
	if l.Name != "Other League" {
		t.Errorf("Name = %q; want Other League", l.Name)
	}
	if l.CurrentRulesVersion != 3 {
		t.Errorf("CurrentRulesVersion = %d; want 3", l.CurrentRulesVersion)
	}
}

func TestLeagueRoundTrip(t *testing.T) {
	l := domain.League{LeagueID: "lrt", Name: "RT League", CurrentRulesVersion: 1}
	back := ItemToLeague(LeagueToItem(l))
	if back.LeagueID != l.LeagueID || back.Name != l.Name || back.CurrentRulesVersion != l.CurrentRulesVersion {
		t.Errorf("round-trip mismatch: %+v", back)
	}
}
