package storage

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
)

type LeagueItem struct {
	PK                  string `dynamodbav:"PK"`
	SK                  string `dynamodbav:"SK"`
	LeagueID            string `dynamodbav:"leagueId"`
	Name                string `dynamodbav:"name"`
	CurrentRulesVersion int    `dynamodbav:"currentRulesVersion"`
}

func LeagueToItem(l domain.League) LeagueItem {
	return LeagueItem{
		PK:                  keys.LeaguePK(l.LeagueID),
		SK:                  keys.LeagueSK(l.LeagueID),
		LeagueID:            l.LeagueID,
		Name:                l.Name,
		CurrentRulesVersion: l.CurrentRulesVersion,
	}
}

func ItemToLeague(it LeagueItem) domain.League {
	return domain.League{
		LeagueID:            it.LeagueID,
		Name:                it.Name,
		CurrentRulesVersion: it.CurrentRulesVersion,
	}
}
