package mappers

import (
	"testing"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
)

func TestPlayerToItem(t *testing.T) {
	p := domain.Player{
		PlayerID: "p1",
		TeamID:   "t1",
		LeagueID: "l1",
		Name:     "Alice",
		Number:   7,
		Position: "pitcher",
	}
	it := PlayerToItem(p)
	if it.PK != keys.TeamPK("t1") {
		t.Errorf("PK = %q; want %q", it.PK, keys.TeamPK("t1"))
	}
	if it.SK != keys.PlayerSK("p1") {
		t.Errorf("SK = %q; want %q", it.SK, keys.PlayerSK("p1"))
	}
	if it.PlayerID != "p1" {
		t.Errorf("PlayerID = %q; want p1", it.PlayerID)
	}
	if it.Name != "Alice" {
		t.Errorf("Name = %q; want Alice", it.Name)
	}
	if it.Number != 7 {
		t.Errorf("Number = %d; want 7", it.Number)
	}
	if it.Position != "pitcher" {
		t.Errorf("Position = %q; want pitcher", it.Position)
	}
}

func TestItemToPlayer(t *testing.T) {
	it := storage.PlayerItem{
		PlayerID: "p2",
		TeamID:   "t2",
		LeagueID: "l2",
		Name:     "Bob",
		Number:   9,
		Position: "catcher",
	}
	p := ItemToPlayer(it)
	if p.PlayerID != "p2" {
		t.Errorf("PlayerID = %q; want p2", p.PlayerID)
	}
	if p.Name != "Bob" {
		t.Errorf("Name = %q; want Bob", p.Name)
	}
	if p.Number != 9 {
		t.Errorf("Number = %d; want 9", p.Number)
	}
}

func TestPlayerRoundTrip(t *testing.T) {
	p := domain.Player{PlayerID: "prt", TeamID: "trt", LeagueID: "lrt", Name: "Carol", Number: 1, Position: "ss"}
	back := ItemToPlayer(PlayerToItem(p))
	if back.PlayerID != p.PlayerID || back.Name != p.Name || back.Number != p.Number {
		t.Errorf("round-trip mismatch: %+v", back)
	}
}
