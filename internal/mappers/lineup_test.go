package mappers

import (
	"testing"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
)

func TestLineupToItem(t *testing.T) {
	l := domain.Lineup{
		GameID:    "g1",
		TeamID:    "t1",
		PlayerIDs: []string{"p1", "p2", "p3"},
	}
	it := LineupToItem(l)
	if it.PK != keys.GamePK("g1") {
		t.Errorf("PK = %q; want %q", it.PK, keys.GamePK("g1"))
	}
	if it.SK != keys.LineupSK("t1") {
		t.Errorf("SK = %q; want %q", it.SK, keys.LineupSK("t1"))
	}
	if it.GameID != "g1" {
		t.Errorf("GameID = %q; want g1", it.GameID)
	}
	if it.TeamID != "t1" {
		t.Errorf("TeamID = %q; want t1", it.TeamID)
	}
	if len(it.PlayerIDs) != 3 || it.PlayerIDs[0] != "p1" {
		t.Errorf("PlayerIDs mismatch: %v", it.PlayerIDs)
	}
}

func TestItemToLineup(t *testing.T) {
	it := storage.LineupItem{
		GameID:    "g2",
		TeamID:    "t2",
		PlayerIDs: []string{"pa", "pb"},
	}
	l := ItemToLineup(it)
	if l.GameID != "g2" {
		t.Errorf("GameID = %q; want g2", l.GameID)
	}
	if l.TeamID != "t2" {
		t.Errorf("TeamID = %q; want t2", l.TeamID)
	}
	if len(l.PlayerIDs) != 2 || l.PlayerIDs[1] != "pb" {
		t.Errorf("PlayerIDs mismatch: %v", l.PlayerIDs)
	}
}

func TestLineupRoundTrip(t *testing.T) {
	l := domain.Lineup{GameID: "grt", TeamID: "trt", PlayerIDs: []string{"px"}}
	back := ItemToLineup(LineupToItem(l))
	if back.GameID != l.GameID || back.TeamID != l.TeamID || len(back.PlayerIDs) != 1 {
		t.Errorf("round-trip mismatch: %+v", back)
	}
}
