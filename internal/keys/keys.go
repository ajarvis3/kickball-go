package keys

import "fmt"

func LeaguePK(leagueID string) string { return fmt.Sprintf("LEAGUE#%s", leagueID) }
func LeagueSK(leagueID string) string { return fmt.Sprintf("LEAGUE#%s", leagueID) }

func TeamPK(teamID string) string { return fmt.Sprintf("TEAM#%s", teamID) }
func TeamSK(teamID string) string { return fmt.Sprintf("TEAM#%s", teamID) }

func GamePK(gameID string) string { return fmt.Sprintf("GAME#%s", gameID) }
func GameSK(gameID string) string { return fmt.Sprintf("GAME#%s", gameID) }

func PlayerSK(playerID string) string { return fmt.Sprintf("PLAYER#%s", playerID) }

func AtBatSK(seq int) string { return fmt.Sprintf("ATBAT#%04d", seq) }

func TeamGameSK(gameID string) string { return fmt.Sprintf("GAME#%s", gameID) }

func GSI2PK(playerID string) string { return fmt.Sprintf("PLAYER#%s", playerID) }
func GSI2SK(gameID string, seq int) string {
	return fmt.Sprintf("GAME#%s#ATBAT#%04d", gameID, seq)
}

func RulesSK(version int) string { return fmt.Sprintf("RULES#%d", version) }

func LineupSK(teamID string) string { return fmt.Sprintf("LINEUP#%s", teamID) }
