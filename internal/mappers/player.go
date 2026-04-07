package mappers

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
)

func PlayerToItem(p domain.Player) storage.PlayerItem {
	return storage.PlayerItem{
		PK:       keys.TeamPK(p.TeamID),
		SK:       keys.PlayerSK(p.PlayerID),
		PlayerID: p.PlayerID,
		TeamID:   p.TeamID,
		LeagueID: p.LeagueID,
		Name:     p.Name,
		Number:   p.Number,
		Position: p.Position,
	}
}

func ItemToPlayer(it storage.PlayerItem) domain.Player {
	return domain.Player{
		PlayerID: it.PlayerID,
		TeamID:   it.TeamID,
		LeagueID: it.LeagueID,
		Name:     it.Name,
		Number:   it.Number,
		Position: it.Position,
	}
}
