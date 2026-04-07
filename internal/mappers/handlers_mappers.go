package mappers

import (
	"github.com/ajarvis3/kickball-go/internal/domain"
	dto "github.com/ajarvis3/kickball-go/internal/handlers/dto"
)

// AtBat mappers
func RecordAtBatRequestToDomain(r dto.RecordAtBatRequest, gameID, leagueID string) domain.AtBat {
	return domain.AtBat{
		GameID:   gameID,
		LeagueID: leagueID,
		TeamID:   r.TeamID,
		PlayerID: r.PlayerID,
		Seq:      1,
		Inning:   0,
		Half:     "",
		Strikes:  r.Strikes,
		Balls:    r.Balls,
		Fouls:    r.Fouls,
		Result:   r.Result,
		RBI:      r.RBI,
		Pitches:  r.Pitches,
	}
}

func AtBatToResponse(a domain.AtBat) dto.AtBatResponse {
	return dto.AtBatResponse{
		GameID:   a.GameID,
		PlayerID: a.PlayerID,
		TeamID:   a.TeamID,
		Seq:      a.Seq,
		Inning:   a.Inning,
		Half:     a.Half,
		Strikes:  a.Strikes,
		Balls:    a.Balls,
		Fouls:    a.Fouls,
		Result:   a.Result,
		RBI:      a.RBI,
		Pitches:  a.Pitches,
	}
}

// Team mappers
func CreateTeamRequestToDomain(r dto.CreateTeamRequest, teamID, leagueID string) domain.Team {
	return domain.Team{TeamID: teamID, LeagueID: leagueID, Name: r.Name}
}

func TeamToResponse(t domain.Team) dto.TeamResponse {
	return dto.TeamResponse{TeamID: t.TeamID, LeagueID: t.LeagueID, Name: t.Name}
}

// Player mappers
func CreatePlayerRequestToDomain(r dto.CreatePlayerRequest, playerID, teamID, leagueID string) domain.Player {
	return domain.Player{PlayerID: playerID, TeamID: teamID, LeagueID: leagueID, Name: r.Name, Number: r.Number, Position: r.Position}
}

func PlayerToResponse(p domain.Player) dto.PlayerResponse {
	return dto.PlayerResponse{PlayerID: p.PlayerID, TeamID: p.TeamID, LeagueID: p.LeagueID, Name: p.Name, Number: p.Number, Position: p.Position}
}

// League mappers
func CreateLeagueRequestToDomain(r dto.CreateLeagueRequest, leagueID string) domain.League {
	return domain.League{LeagueID: leagueID, Name: r.Name, CurrentRulesVersion: 1}
}

func LeagueToResponse(l domain.League) dto.LeagueResponse {
	return dto.LeagueResponse{LeagueID: l.LeagueID, Name: l.Name, CurrentRulesVersion: l.CurrentRulesVersion}
}

// LeagueRules mappers
func CreateLeagueRulesRequestToDomain(r dto.CreateLeagueRulesRequest, leagueID string, version int) domain.LeagueRules {
	return domain.LeagueRules{
		LeagueID:               leagueID,
		RulesVersion:           version,
		MaxStrikes:             r.MaxStrikes,
		MaxBalls:               r.MaxBalls,
		MaxFouls:               r.MaxFouls,
		MaxInnings:             r.MaxInnings,
		MercyRunsPerInning:     r.MercyRunsPerInning,
		MercyAppliesLastInning: r.MercyAppliesLastInning,
		GameMercyRuns:          r.GameMercyRuns,
		Metadata:               r.Metadata,
	}
}

func LeagueRulesToResponse(r domain.LeagueRules) dto.LeagueRulesResponse {
	return dto.LeagueRulesResponse{
		LeagueID:               r.LeagueID,
		RulesVersion:           r.RulesVersion,
		MaxStrikes:             r.MaxStrikes,
		MaxBalls:               r.MaxBalls,
		MaxFouls:               r.MaxFouls,
		MaxInnings:             r.MaxInnings,
		MercyRunsPerInning:     r.MercyRunsPerInning,
		MercyAppliesLastInning: r.MercyAppliesLastInning,
		GameMercyRuns:          r.GameMercyRuns,
		Metadata:               r.Metadata,
	}
}

// Lineup mappers
func SetLineupRequestToDomain(r dto.SetLineupRequest, gameID, teamID string) domain.Lineup {
	return domain.Lineup{GameID: gameID, TeamID: teamID, PlayerIDs: r.PlayerIDs}
}

func LineupToResponse(l domain.Lineup) dto.LineupResponse {
	return dto.LineupResponse{GameID: l.GameID, TeamID: l.TeamID, PlayerIDs: l.PlayerIDs}
}

// Game mappers
func GameToResponse(g domain.Game) dto.GameResponse {
	return dto.GameResponse{
		GameID:       g.GameID,
		LeagueID:     g.LeagueID,
		HomeTeamID:   g.HomeTeamID,
		AwayTeamID:   g.AwayTeamID,
		RulesVersion: g.RulesVersion,
		State: dto.GameStateDTO{
			Inning:    g.State.Inning,
			Half:      g.State.Half,
			Outs:      g.State.Outs,
			HomeScore: g.State.HomeScore,
			AwayScore: g.State.AwayScore,
		},
	}
}

// CreateGameRequestToDomain converts a CreateGameRequest to a domain.Game with
// the provided gameID and leagueID. State is left empty for the handler to
// initialize after loading league rules.
func CreateGameRequestToDomain(r dto.CreateGameRequest, gameID, leagueID string) domain.Game {
	return domain.Game{
		GameID:       gameID,
		LeagueID:     leagueID,
		RulesVersion: 1,
		HomeTeamID:   r.HomeTeamID,
		AwayTeamID:   r.AwayTeamID,
		State:        domain.GameState{},
	}
}

// TeamGame response
func TeamGameToResponse(t domain.TeamGame) dto.TeamGameResponse {
	return dto.TeamGameResponse{GameID: t.GameID, TeamID: t.TeamID, OpponentID: t.OpponentID, LeagueID: t.LeagueID, Date: t.Date}
}
