package storage

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
)

type TeamGameItem struct {
	PK         string `dynamodbav:"PK"`
	SK         string `dynamodbav:"SK"`
	GameID     string `dynamodbav:"GameID"`
	TeamID     string `dynamodbav:"TeamID"`
	OpponentID string `dynamodbav:"OpponentID"`
	LeagueID   string `dynamodbav:"LeagueID"`
	Date       string `dynamodbav:"date"`
}

func TeamGameToItem(t domain.TeamGame) TeamGameItem {
	return TeamGameItem{
		PK:         keys.TeamPK(t.TeamID),
		SK:         keys.TeamGameSK(t.GameID),
		GameID:     t.GameID,
		TeamID:     t.TeamID,
		OpponentID: t.OpponentID,
		LeagueID:   t.LeagueID,
		Date:       t.Date,
	}
}

func ItemToTeamGame(it TeamGameItem) domain.TeamGame {
	return domain.TeamGame{
		GameID:     it.GameID,
		TeamID:     it.TeamID,
		OpponentID: it.OpponentID,
		LeagueID:   it.LeagueID,
		Date:       it.Date,
	}
}
