package mappers

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
)

func GameToItem(g domain.Game) storage.GameItem {
	return storage.GameItem{
		PK:              keys.GamePK(g.GameID),
		SK:              keys.GameSK(g.GameID),
		GameID:          g.GameID,
		LeagueID:        g.LeagueID,
		RulesVersion:    g.RulesVersion,
		HomeTeamID:      g.HomeTeamID,
		AwayTeamID:      g.AwayTeamID,
		State:           g.State,
		GSILeagueGamePK: keys.LeaguePK(g.LeagueID),
		GSILeagueGameSK: keys.TeamGameSK(g.GameID),
	}
}

func ItemToGame(it storage.GameItem) domain.Game {
	return domain.Game{
		GameID:       it.GameID,
		LeagueID:     it.LeagueID,
		RulesVersion: it.RulesVersion,
		HomeTeamID:   it.HomeTeamID,
		AwayTeamID:   it.AwayTeamID,
		State:        it.State,
	}
}
