package mappers

import (
	"testing"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
)

func TestGameToItem(t *testing.T) {
	g := domain.Game{
		GameID:       "g1",
		LeagueID:     "l1",
		RulesVersion: 1,
		HomeTeamID:   "ht1",
		AwayTeamID:   "at1",
		State: domain.GameState{
			Inning:    1,
			Half:      "top",
			Outs:      0,
			HomeScore: 0,
			AwayScore: 0,
		},
	}
	it := GameToItem(g)
	if it.PK != keys.GamePK("g1") {
		t.Errorf("PK = %q; want %q", it.PK, keys.GamePK("g1"))
	}
	if it.SK != keys.GameSK("g1") {
		t.Errorf("SK = %q; want %q", it.SK, keys.GameSK("g1"))
	}
	if it.GameID != "g1" {
		t.Errorf("GameID = %q; want g1", it.GameID)
	}
	if it.LeagueID != "l1" {
		t.Errorf("LeagueID = %q; want l1", it.LeagueID)
	}
	if it.HomeTeamID != "ht1" {
		t.Errorf("HomeTeamID = %q; want ht1", it.HomeTeamID)
	}
	if it.AwayTeamID != "at1" {
		t.Errorf("AwayTeamID = %q; want at1", it.AwayTeamID)
	}
	if it.GSILeagueGamePK != keys.LeaguePK("l1") {
		t.Errorf("GSILeagueGamePK = %q; want %q", it.GSILeagueGamePK, keys.LeaguePK("l1"))
	}
	if it.GSILeagueGameSK != keys.TeamGameSK("g1") {
		t.Errorf("GSILeagueGameSK = %q; want %q", it.GSILeagueGameSK, keys.TeamGameSK("g1"))
	}
}

func TestItemToGame(t *testing.T) {
	it := storage.GameItem{
		GameID:       "g2",
		LeagueID:     "l2",
		RulesVersion: 2,
		HomeTeamID:   "ht2",
		AwayTeamID:   "at2",
		State: domain.GameState{
			Inning:    3,
			Half:      "bottom",
			Outs:      2,
			HomeScore: 4,
			AwayScore: 2,
		},
	}
	g := ItemToGame(it)
	if g.GameID != "g2" {
		t.Errorf("GameID = %q; want g2", g.GameID)
	}
	if g.LeagueID != "l2" {
		t.Errorf("LeagueID = %q; want l2", g.LeagueID)
	}
	if g.RulesVersion != 2 {
		t.Errorf("RulesVersion = %d; want 2", g.RulesVersion)
	}
	if g.State.Inning != 3 {
		t.Errorf("State.Inning = %d; want 3", g.State.Inning)
	}
	if g.State.HomeScore != 4 {
		t.Errorf("State.HomeScore = %d; want 4", g.State.HomeScore)
	}
}

func TestGameRoundTrip(t *testing.T) {
	g := domain.Game{
		GameID:       "grt",
		LeagueID:     "lrt",
		RulesVersion: 1,
		HomeTeamID:   "hrt",
		AwayTeamID:   "art",
		State:        domain.GameState{Inning: 2, Half: "top", Outs: 1},
	}
	back := ItemToGame(GameToItem(g))
	if back.GameID != g.GameID || back.State.Inning != g.State.Inning {
		t.Errorf("round-trip mismatch: %+v", back)
	}
}
