package mappers

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
)

func LeagueToItem(l domain.League) storage.LeagueItem {
	return storage.LeagueItem{
		PK:                  keys.LeaguePK(l.LeagueID),
		SK:                  keys.LeagueSK(l.LeagueID),
		LeagueID:            l.LeagueID,
		Name:                l.Name,
		CurrentRulesVersion: l.CurrentRulesVersion,
	}
}

func ItemToLeague(it storage.LeagueItem) domain.League {
	return domain.League{
		LeagueID:            it.LeagueID,
		Name:                it.Name,
		CurrentRulesVersion: it.CurrentRulesVersion,
	}
}
