package keys

import (
	"fmt"
	"testing"
)

func TestLeaguePK(t *testing.T) {
	got := LeaguePK("abc")
	want := fmt.Sprintf(leaguePrefixFmt, "abc")
	if got != want {
		t.Errorf("LeaguePK = %q; want %q", got, want)
	}
}

func TestLeagueSK(t *testing.T) {
	got := LeagueSK("abc")
	want := fmt.Sprintf(leaguePrefixFmt, "abc")
	if got != want {
		t.Errorf("LeagueSK = %q; want %q", got, want)
	}
}

func TestTeamPK(t *testing.T) {
	got := TeamPK("t1")
	want := fmt.Sprintf(teamPrefixFmt, "t1")
	if got != want {
		t.Errorf("TeamPK = %q; want %q", got, want)
	}
}

func TestTeamSK(t *testing.T) {
	got := TeamSK("t1")
	want := fmt.Sprintf(teamPrefixFmt, "t1")
	if got != want {
		t.Errorf("TeamSK = %q; want %q", got, want)
	}
}

func TestGamePK(t *testing.T) {
	got := GamePK("g1")
	want := fmt.Sprintf(gamePrefixFmt, "g1")
	if got != want {
		t.Errorf("GamePK = %q; want %q", got, want)
	}
}

func TestGameSK(t *testing.T) {
	got := GameSK("g1")
	want := fmt.Sprintf(gamePrefixFmt, "g1")
	if got != want {
		t.Errorf("GameSK = %q; want %q", got, want)
	}
}

func TestPlayerSK(t *testing.T) {
	got := PlayerSK("p1")
	want := fmt.Sprintf(playerPrefixFmt, "p1")
	if got != want {
		t.Errorf("PlayerSK = %q; want %q", got, want)
	}
}

func TestAtBatSK(t *testing.T) {
	got := AtBatSK(3)
	want := fmt.Sprintf(atBatPrefixFmt, 3)
	if got != want {
		t.Errorf("AtBatSK = %q; want %q", got, want)
	}
	// verify zero-padding
	if got != "ATBAT#0003" {
		t.Errorf("AtBatSK zero-padding: got %q; want %q", got, "ATBAT#0003")
	}
}

func TestTeamGameSK(t *testing.T) {
	got := TeamGameSK("g1")
	want := fmt.Sprintf(teamGamePrefixFmt, "g1")
	if got != want {
		t.Errorf("TeamGameSK = %q; want %q", got, want)
	}
}

func TestGSI2PK(t *testing.T) {
	got := GSI2PK("p1")
	want := fmt.Sprintf(gsi2PKFmt, "p1")
	if got != want {
		t.Errorf("GSI2PK = %q; want %q", got, want)
	}
}

func TestGSI2SK(t *testing.T) {
	got := GSI2SK("g1", 5)
	want := fmt.Sprintf(gsi2SKFmt, "g1", 5)
	if got != want {
		t.Errorf("GSI2SK = %q; want %q", got, want)
	}
}

func TestRulesSK(t *testing.T) {
	got := RulesSK(2)
	want := fmt.Sprintf(rulesPrefixFmt, 2)
	if got != want {
		t.Errorf("RulesSK = %q; want %q", got, want)
	}
}

func TestLineupSK(t *testing.T) {
	got := LineupSK("t1")
	want := fmt.Sprintf(lineupPrefixFmt, "t1")
	if got != want {
		t.Errorf("LineupSK = %q; want %q", got, want)
	}
}

func TestLeaguePKFormat(t *testing.T) {
	if LeaguePK("x") != "LEAGUE#x" {
		t.Errorf("LeaguePK format mismatch")
	}
}

func TestTeamPKFormat(t *testing.T) {
	if TeamPK("x") != "TEAM#x" {
		t.Errorf("TeamPK format mismatch")
	}
}

func TestGamePKFormat(t *testing.T) {
	if GamePK("x") != "GAME#x" {
		t.Errorf("GamePK format mismatch")
	}
}

func TestPlayerSKFormat(t *testing.T) {
	if PlayerSK("x") != "PLAYER#x" {
		t.Errorf("PlayerSK format mismatch")
	}
}

func TestGSI2PKFormat(t *testing.T) {
	if GSI2PK("x") != "PLAYER#x" {
		t.Errorf("GSI2PK format mismatch")
	}
}

func TestGSI2SKFormat(t *testing.T) {
	if GSI2SK("g", 1) != "GAME#g#ATBAT#0001" {
		t.Errorf("GSI2SK format mismatch: got %q", GSI2SK("g", 1))
	}
}
