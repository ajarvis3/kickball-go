package storage

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
)

type TeamItem struct {
	PK       string `dynamodbav:"PK"`
	SK       string `dynamodbav:"SK"`
	TeamID   string `dynamodbav:"TeamID"`
	LeagueID string `dynamodbav:"LeagueID"`
	Name     string `dynamodbav:"name"`
}

func TeamToItem(t domain.Team) TeamItem {
	return TeamItem{
		PK:       keys.LeaguePK(t.LeagueID),
		SK:       keys.TeamSK(t.TeamID),
		TeamID:   t.TeamID,
		LeagueID: t.LeagueID,
		Name:     t.Name,
	}
}

func ItemToTeam(it TeamItem) domain.Team {
	return domain.Team{
		TeamID:   it.TeamID,
		LeagueID: it.LeagueID,
		Name:     it.Name,
	}
}
