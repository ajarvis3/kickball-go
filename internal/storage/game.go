package storage

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
)

type GameItem struct {
	PK           string           `dynamodbav:"PK"`
	SK           string           `dynamodbav:"SK"`
	GameID       string           `dynamodbav:"GameID"`
	LeagueID     string           `dynamodbav:"LeagueID"`
	RulesVersion int              `dynamodbav:"rulesVersion"`
	HomeTeamID   string           `dynamodbav:"homeTeamID"`
	AwayTeamID   string           `dynamodbav:"awayTeamID"`
	State        domain.GameState `dynamodbav:"state"`

	GSILeagueGamePK string `dynamodbav:"GSILeagueGamePK"`
	GSILeagueGameSK string `dynamodbav:"GSILeagueGameSK"`
}

func GameToItem(g domain.Game) GameItem {
	return GameItem{
		PK:              keys.GamePK(g.GameID),
		SK:              keys.GameSK(g.GameID),
		GameID:          g.GameID,
		LeagueID:        g.LeagueID,
		RulesVersion:    g.RulesVersion,
		HomeTeamID:      g.HomeTeamID,
		AwayTeamID:      g.AwayTeamID,
		State:           g.State,
		GSILeagueGamePK: keys.LeaguePK(g.LeagueID),
		GSILeagueGameSK: keys.TeamGameSK(g.GameID),
	}
}

func ItemToGame(it GameItem) domain.Game {
	return domain.Game{
		GameID:       it.GameID,
		LeagueID:     it.LeagueID,
		RulesVersion: it.RulesVersion,
		HomeTeamID:   it.HomeTeamID,
		AwayTeamID:   it.AwayTeamID,
		State:        it.State,
	}
}
