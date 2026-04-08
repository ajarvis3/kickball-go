package mappers

import (
	"testing"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
)

func TestAtbatToItem(t *testing.T) {
	a := domain.AtBat{
		GameID:   "g1",
		LeagueID: "l1",
		TeamID:   "t1",
		PlayerID: "p1",
		Seq:      3,
		Inning:   2,
		Half:     "top",
		Strikes:  1,
		Balls:    2,
		Fouls:    0,
		Result:   "single",
		RBI:      1,
		Pitches:  5,
	}
	it := AtbatToItem(a)
	if it.PK != keys.GamePK("g1") {
		t.Errorf("PK = %q; want %q", it.PK, keys.GamePK("g1"))
	}
	if it.SK != keys.AtBatSK(3) {
		t.Errorf("SK = %q; want %q", it.SK, keys.AtBatSK(3))
	}
	if it.GameID != "g1" {
		t.Errorf("GameID = %q; want g1", it.GameID)
	}
	if it.LeagueID != "l1" {
		t.Errorf("LeagueID = %q; want l1", it.LeagueID)
	}
	if it.TeamID != "t1" {
		t.Errorf("TeamID = %q; want t1", it.TeamID)
	}
	if it.PlayerID != "p1" {
		t.Errorf("PlayerID = %q; want p1", it.PlayerID)
	}
	if it.Seq != 3 {
		t.Errorf("Seq = %d; want 3", it.Seq)
	}
	if it.GSIPlayerAtBatPK != keys.GSI2PK("p1") {
		t.Errorf("GSIPlayerAtBatPK = %q; want %q", it.GSIPlayerAtBatPK, keys.GSI2PK("p1"))
	}
	if it.GSIPlayerAtBatSK != keys.GSI2SK("g1", 3) {
		t.Errorf("GSIPlayerAtBatSK = %q; want %q", it.GSIPlayerAtBatSK, keys.GSI2SK("g1", 3))
	}
	if it.Result != "single" {
		t.Errorf("Result = %q; want single", it.Result)
	}
	if it.RBI != 1 {
		t.Errorf("RBI = %d; want 1", it.RBI)
	}
}

func TestItemToAtbat(t *testing.T) {
	it := storage.AtbatItem{
		GameID:   "g2",
		LeagueID: "l2",
		TeamID:   "t2",
		PlayerID: "p2",
		Seq:      5,
		Inning:   3,
		Half:     "bottom",
		Strikes:  2,
		Balls:    1,
		Fouls:    3,
		Result:   "out",
		RBI:      0,
		Pitches:  7,
	}
	a := ItemToAtbat(it)
	if a.GameID != "g2" {
		t.Errorf("GameID = %q; want g2", a.GameID)
	}
	if a.PlayerID != "p2" {
		t.Errorf("PlayerID = %q; want p2", a.PlayerID)
	}
	if a.Seq != 5 {
		t.Errorf("Seq = %d; want 5", a.Seq)
	}
	if a.Half != "bottom" {
		t.Errorf("Half = %q; want bottom", a.Half)
	}
	if a.Result != "out" {
		t.Errorf("Result = %q; want out", a.Result)
	}
	if a.Pitches != 7 {
		t.Errorf("Pitches = %d; want 7", a.Pitches)
	}
}

func TestAtbatRoundTrip(t *testing.T) {
	a := domain.AtBat{
		GameID:   "gx",
		LeagueID: "lx",
		TeamID:   "tx",
		PlayerID: "px",
		Seq:      9,
		Inning:   4,
		Half:     "top",
		Strikes:  3,
		Balls:    3,
		Fouls:    2,
		Result:   "homerun",
		RBI:      3,
		Pitches:  10,
	}
	back := ItemToAtbat(AtbatToItem(a))
	if back.GameID != a.GameID || back.Seq != a.Seq || back.Result != a.Result {
		t.Errorf("round-trip mismatch: %+v", back)
	}
}
