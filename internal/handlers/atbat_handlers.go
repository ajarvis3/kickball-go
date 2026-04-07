package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/ajarvis3/kickball-go/internal/db"
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/handlers/dto"
	"github.com/ajarvis3/kickball-go/internal/services"
	"github.com/ajarvis3/kickball-go/pkg/responses"
)

type AtBatHandlers struct {
	AtBats db.AtBatRepository
	Games  db.GameRepository
	Rules  db.LeagueRulesRepository
	Engine *services.GameEngine
}

func NewAtBatHandlers(atbats db.AtBatRepository, games db.GameRepository, rules db.LeagueRulesRepository, engine *services.GameEngine) *AtBatHandlers {
	return &AtBatHandlers{AtBats: atbats, Games: games, Rules: rules, Engine: engine}
}

func (h *AtBatHandlers) RecordAtBat(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body dto.RecordAtBatRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}), nil
	}
	gameID := req.PathParameters["gameId"]
	leagueID := req.PathParameters["leagueId"]
	if gameID == "" || body.PlayerID == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "gameId and playerId are required"}), nil
	}
	atbat := domain.AtBat{GameID: gameID, LeagueID: leagueID, TeamID: body.TeamID, PlayerID: body.PlayerID, Seq: 1, Inning: 0, Half: "", Strikes: body.Strikes, Balls: body.Balls, Fouls: body.Fouls, Result: body.Result, RBI: body.RBI, Pitches: body.Pitches}

	// Load current game
	g, err := h.Games.GetGame(ctx, gameID)
	if err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}
	if g == nil {
		return responses.JsonResponse(http.StatusNotFound, map[string]string{"error": "game not found"}), nil
	}

	// Load league rules for the game's rules version
	lr, err := h.Rules.GetLeagueRules(ctx, leagueID, g.RulesVersion)
	if err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}
	if lr == nil {
		return responses.JsonResponse(http.StatusNotFound, map[string]string{"error": "league rules not found"}), nil
	}

	// Apply at-bat to game state
	updatedGame, err := h.Engine.ApplyAtBat(*g, *lr, atbat)
	if err != nil {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}), nil
	}

	// Persist at-bat and updated game transactionally
	if err := h.AtBats.PutAtBatAndUpdateGame(ctx, atbat, updatedGame); err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}

	resp := dto.AtBatResponse{GameID: atbat.GameID, PlayerID: atbat.PlayerID, TeamID: atbat.TeamID, Seq: atbat.Seq, Inning: atbat.Inning, Half: atbat.Half, Strikes: atbat.Strikes, Balls: atbat.Balls, Fouls: atbat.Fouls, Result: atbat.Result, RBI: atbat.RBI, Pitches: atbat.Pitches}
	return responses.JsonResponse(http.StatusCreated, resp), nil
}

func (h *AtBatHandlers) GetAtBats(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Support query params: gameId or playerId
	gameID := req.QueryStringParameters["gameId"]
	playerID := req.QueryStringParameters["playerId"]
	if gameID == "" && playerID == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "gameId or playerId query parameter is required"}), nil
	}

	var (
		atbats []domain.AtBat
		err    error
	)

	if gameID != "" {
		atbats, err = h.AtBats.ListAtBatsByGame(ctx, gameID)
	} else {
		atbats, err = h.AtBats.ListAtBatsByPlayer(ctx, playerID)
	}
	if err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}

	var resp []dto.AtBatResponse
	for _, a := range atbats {
		resp = append(resp, dto.AtBatResponse{GameID: a.GameID, PlayerID: a.PlayerID, TeamID: a.TeamID, Seq: a.Seq, Inning: a.Inning, Half: a.Half, Strikes: a.Strikes, Balls: a.Balls, Fouls: a.Fouls, Result: a.Result, RBI: a.RBI, Pitches: a.Pitches})
	}

	return responses.JsonResponse(http.StatusOK, resp), nil
}
