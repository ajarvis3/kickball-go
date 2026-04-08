package mappers

import (
	"testing"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
)

func TestTeamGameToItem(t *testing.T) {
	tg := domain.TeamGame{
		GameID:     "g1",
		TeamID:     "t1",
		OpponentID: "t2",
		LeagueID:   "l1",
		Date:       "2024-01-01",
	}
	it := TeamGameToItem(tg)
	if it.PK != keys.TeamPK("t1") {
		t.Errorf("PK = %q; want %q", it.PK, keys.TeamPK("t1"))
	}
	if it.SK != keys.TeamGameSK("g1") {
		t.Errorf("SK = %q; want %q", it.SK, keys.TeamGameSK("g1"))
	}
	if it.GameID != "g1" {
		t.Errorf("GameID = %q; want g1", it.GameID)
	}
	if it.TeamID != "t1" {
		t.Errorf("TeamID = %q; want t1", it.TeamID)
	}
	if it.OpponentID != "t2" {
		t.Errorf("OpponentID = %q; want t2", it.OpponentID)
	}
	if it.Date != "2024-01-01" {
		t.Errorf("Date = %q; want 2024-01-01", it.Date)
	}
}

func TestItemToTeamGame(t *testing.T) {
	it := storage.TeamGameItem{
		GameID:     "g2",
		TeamID:     "t3",
		OpponentID: "t4",
		LeagueID:   "l2",
		Date:       "2024-06-15",
	}
	tg := ItemToTeamGame(it)
	if tg.GameID != "g2" {
		t.Errorf("GameID = %q; want g2", tg.GameID)
	}
	if tg.OpponentID != "t4" {
		t.Errorf("OpponentID = %q; want t4", tg.OpponentID)
	}
	if tg.Date != "2024-06-15" {
		t.Errorf("Date = %q; want 2024-06-15", tg.Date)
	}
}

func TestTeamGameRoundTrip(t *testing.T) {
	tg := domain.TeamGame{GameID: "grt", TeamID: "trt", OpponentID: "opp", LeagueID: "lrt", Date: "2024-01-01"}
	back := ItemToTeamGame(TeamGameToItem(tg))
	if back.GameID != tg.GameID || back.OpponentID != tg.OpponentID {
		t.Errorf("round-trip mismatch: %+v", back)
	}
}
