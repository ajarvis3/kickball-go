package mappers

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
)

func TeamGameToItem(t domain.TeamGame) storage.TeamGameItem {
	return storage.TeamGameItem{
		PK:         keys.TeamPK(t.TeamID),
		SK:         keys.TeamGameSK(t.GameID),
		GameID:     t.GameID,
		TeamID:     t.TeamID,
		OpponentID: t.OpponentID,
		LeagueID:   t.LeagueID,
		Date:       t.Date,
	}
}

func ItemToTeamGame(it storage.TeamGameItem) domain.TeamGame {
	return domain.TeamGame{
		GameID:     it.GameID,
		TeamID:     it.TeamID,
		OpponentID: it.OpponentID,
		LeagueID:   it.LeagueID,
		Date:       it.Date,
	}
}
