package mappers

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
)

func TeamToItem(t domain.Team) storage.TeamItem {
	return storage.TeamItem{
		PK:       keys.LeaguePK(t.LeagueID),
		SK:       keys.TeamSK(t.TeamID),
		TeamID:   t.TeamID,
		LeagueID: t.LeagueID,
		Name:     t.Name,
	}
}

func ItemToTeam(it storage.TeamItem) domain.Team {
	return domain.Team{
		TeamID:   it.TeamID,
		LeagueID: it.LeagueID,
		Name:     it.Name,
	}
}
