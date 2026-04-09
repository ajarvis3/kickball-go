package mappers

import (
	"github.com/ajarvis3/kickball-go/internal/data/domain"
	"github.com/ajarvis3/kickball-go/internal/data/storage"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/types"
)

func AtbatToItem(a domain.AtBat) storage.AtbatItem {
	return storage.AtbatItem{
		PK:               keys.GamePK(a.GameID),
		SK:               keys.AtBatSK(a.Seq),
		GameID:           a.GameID,
		LeagueID:         a.LeagueID,
		TeamID:           a.TeamID,
		PlayerID:         a.PlayerID,
		Seq:              a.Seq,
		Inning:           a.Inning,
		Half:             string(a.Half),
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

func ItemToAtbat(it storage.AtbatItem) domain.AtBat {
	return domain.AtBat{
		GameID:   it.GameID,
		LeagueID: it.LeagueID,
		TeamID:   it.TeamID,
		PlayerID: it.PlayerID,
		Seq:      it.Seq,
		Inning:   it.Inning,
		Half:     types.Half(it.Half),
		Strikes:  it.Strikes,
		Balls:    it.Balls,
		Fouls:    it.Fouls,
		Result:   it.Result,
		RBI:      it.RBI,
		Pitches:  it.Pitches,
	}
}
