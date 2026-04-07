package storage

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
)

type PlayerItem struct {
	PK       string `dynamodbav:"PK"`
	SK       string `dynamodbav:"SK"`
	PlayerID string `dynamodbav:"PlayerID"`
	TeamID   string `dynamodbav:"TeamID"`
	LeagueID string `dynamodbav:"LeagueID"`
	Name     string `dynamodbav:"name"`
	Number   int    `dynamodbav:"number"`
	Position string `dynamodbav:"position"`
}

func PlayerToItem(p domain.Player) PlayerItem {
	return PlayerItem{
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

func ItemToPlayer(it PlayerItem) domain.Player {
	return domain.Player{
		PlayerID: it.PlayerID,
		TeamID:   it.TeamID,
		LeagueID: it.LeagueID,
		Name:     it.Name,
		Number:   it.Number,
		Position: it.Position,
	}
}
