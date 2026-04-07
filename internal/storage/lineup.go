package storage

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
)

type LineupItem struct {
	PK        string   `dynamodbav:"PK"`
	SK        string   `dynamodbav:"SK"`
	GameID    string   `dynamodbav:"GameID"`
	TeamID    string   `dynamodbav:"TeamID"`
	PlayerIDs []string `dynamodbav:"playerIds"`
}

func LineupToItem(l domain.Lineup) LineupItem {
	return LineupItem{
		PK:        keys.GamePK(l.GameID),
		SK:        keys.LineupSK(l.TeamID),
		GameID:    l.GameID,
		TeamID:    l.TeamID,
		PlayerIDs: l.PlayerIDs,
	}
}

func ItemToLineup(it LineupItem) domain.Lineup {
	return domain.Lineup{
		GameID:    it.GameID,
		TeamID:    it.TeamID,
		PlayerIDs: it.PlayerIDs,
	}
}
