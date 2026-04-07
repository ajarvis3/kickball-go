package keys

import "fmt"

const (
	leaguePrefixFmt   = "LEAGUE#%s"
	teamPrefixFmt     = "TEAM#%s"
	gamePrefixFmt     = "GAME#%s"
	playerPrefixFmt   = "PLAYER#%s"
	atBatPrefixFmt    = "ATBAT#%04d"
	teamGamePrefixFmt = "GAME#%s"
	gsi2PKFmt         = "PLAYER#%s"
	gsi2SKFmt         = "GAME#%s#ATBAT#%04d"
	rulesPrefixFmt    = "RULES#%d"
	lineupPrefixFmt   = "LINEUP#%s"
)

func LeaguePK(leagueID string) string { return fmt.Sprintf(leaguePrefixFmt, leagueID) }
func LeagueSK(leagueID string) string { return fmt.Sprintf(leaguePrefixFmt, leagueID) }

func TeamPK(teamID string) string { return fmt.Sprintf(teamPrefixFmt, teamID) }
func TeamSK(teamID string) string { return fmt.Sprintf(teamPrefixFmt, teamID) }

func GamePK(gameID string) string { return fmt.Sprintf(gamePrefixFmt, gameID) }
func GameSK(gameID string) string { return fmt.Sprintf(gamePrefixFmt, gameID) }

func PlayerSK(playerID string) string { return fmt.Sprintf(playerPrefixFmt, playerID) }

func AtBatSK(seq int) string { return fmt.Sprintf(atBatPrefixFmt, seq) }

func TeamGameSK(gameID string) string { return fmt.Sprintf(teamGamePrefixFmt, gameID) }

func GSI2PK(playerID string) string        { return fmt.Sprintf(gsi2PKFmt, playerID) }
func GSI2SK(gameID string, seq int) string { return fmt.Sprintf(gsi2SKFmt, gameID, seq) }

func RulesSK(version int) string { return fmt.Sprintf(rulesPrefixFmt, version) }

func LineupSK(teamID string) string { return fmt.Sprintf(lineupPrefixFmt, teamID) }
