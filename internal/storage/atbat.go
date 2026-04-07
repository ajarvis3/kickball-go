package storage

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
)

type AtbatItem struct {
	PK               string `dynamodbav:"PK"`
	SK               string `dynamodbav:"SK"`
	GameID           string `dynamodbav:"GameID"`
	LeagueID         string `dynamodbav:"LeagueID"`
	TeamID           string `dynamodbav:"TeamID"`
	PlayerID         string `dynamodbav:"PlayerID"`
	Seq              int    `dynamodbav:"seq"`
	Inning           int    `dynamodbav:"inning"`
	Half             string `dynamodbav:"half"`
	Strikes          int    `dynamodbav:"strikes"`
	Balls            int    `dynamodbav:"balls"`
	Fouls            int    `dynamodbav:"fouls"`
	Result           string `dynamodbav:"result"`
	RBI              int    `dynamodbav:"rbi"`
	Pitches          int    `dynamodbav:"pitches"`
	GSIPlayerAtBatPK string `dynamodbav:"GSIPlayerAtBatPK"`
	GSIPlayerAtBatSK string `dynamodbav:"GSIPlayerAtBatSK"`
}

func AtbatToItem(a domain.AtBat) AtbatItem {
	return AtbatItem{
		PK:               keys.GamePK(a.GameID),
		SK:               keys.AtBatSK(a.Seq),
		GameID:           a.GameID,
		LeagueID:         a.LeagueID,
		TeamID:           a.TeamID,
		PlayerID:         a.PlayerID,
		Seq:              a.Seq,
		Inning:           a.Inning,
		Half:             a.Half,
		Strikes:          a.Strikes,
		Balls:            a.Balls,
		Fouls:            a.Fouls,
		Result:           a.Result,
		RBI:              a.RBI,
		Pitches:          a.Pitches,
		GSIPlayerAtBatPK: keys.GSI2PK(a.PlayerID),
		GSIPlayerAtBatSK: keys.GSI2SK(a.GameID, a.Seq),
	}
}

func ItemToAtbat(it AtbatItem) domain.AtBat {
	return domain.AtBat{
		GameID:   it.GameID,
		LeagueID: it.LeagueID,
		TeamID:   it.TeamID,
		PlayerID: it.PlayerID,
		Seq:      it.Seq,
		Inning:   it.Inning,
		Half:     it.Half,
		Strikes:  it.Strikes,
		Balls:    it.Balls,
		Fouls:    it.Fouls,
		Result:   it.Result,
		RBI:      it.RBI,
		Pitches:  it.Pitches,
	}
}
