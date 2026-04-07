package mappers

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
)

func LineupToItem(l domain.Lineup) storage.LineupItem {
	return storage.LineupItem{
		PK:        keys.GamePK(l.GameID),
		SK:        keys.LineupSK(l.TeamID),
		GameID:    l.GameID,
		TeamID:    l.TeamID,
		PlayerIDs: l.PlayerIDs,
	}
}

func ItemToLineup(it storage.LineupItem) domain.Lineup {
	return domain.Lineup{
		GameID:    it.GameID,
		TeamID:    it.TeamID,
		PlayerIDs: it.PlayerIDs,
	}
}
