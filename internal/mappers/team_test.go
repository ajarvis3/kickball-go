package mappers

import (
	"testing"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
)

func TestTeamToItem(t *testing.T) {
	tm := domain.Team{
		TeamID:   "t1",
		LeagueID: "l1",
		Name:     "Red Sox",
	}
	it := TeamToItem(tm)
	if it.PK != keys.LeaguePK("l1") {
		t.Errorf("PK = %q; want %q", it.PK, keys.LeaguePK("l1"))
	}
	if it.SK != keys.TeamSK("t1") {
		t.Errorf("SK = %q; want %q", it.SK, keys.TeamSK("t1"))
	}
	if it.TeamID != "t1" {
		t.Errorf("TeamID = %q; want t1", it.TeamID)
	}
	if it.LeagueID != "l1" {
		t.Errorf("LeagueID = %q; want l1", it.LeagueID)
	}
	if it.Name != "Red Sox" {
		t.Errorf("Name = %q; want Red Sox", it.Name)
	}
}

func TestItemToTeam(t *testing.T) {
	it := storage.TeamItem{
		TeamID:   "t2",
		LeagueID: "l2",
		Name:     "Blue Jays",
	}
	tm := ItemToTeam(it)
	if tm.TeamID != "t2" {
		t.Errorf("TeamID = %q; want t2", tm.TeamID)
	}
	if tm.Name != "Blue Jays" {
		t.Errorf("Name = %q; want Blue Jays", tm.Name)
	}
}

func TestTeamRoundTrip(t *testing.T) {
	tm := domain.Team{TeamID: "trt", LeagueID: "lrt", Name: "Cubs"}
	back := ItemToTeam(TeamToItem(tm))
	if back.TeamID != tm.TeamID || back.Name != tm.Name {
		t.Errorf("round-trip mismatch: %+v", back)
	}
}
