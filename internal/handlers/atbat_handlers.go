package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/ajarvis3/kickball-go/internal/data/domain"
	dto "github.com/ajarvis3/kickball-go/internal/data/dto"
	"github.com/ajarvis3/kickball-go/internal/db"
	"github.com/ajarvis3/kickball-go/internal/mappers"
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
	// Accept gameId and leagueId from request body instead of path parameters
	gameID := body.GameID
	leagueID := body.LeagueID
	if gameID == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "gameId is required"}), nil
	}
	atbat := mappers.RecordAtBatRequestToDomain(body, gameID, leagueID)

	// Load current game
	g, resp := fetchResource(func() (*domain.Game, error) { return h.Games.GetGame(ctx, gameID) }, "game not found")
	if resp != nil {
		return *resp, nil
	}

	// Load league rules for the game's rules version
	lr, resp := fetchResource(func() (*domain.LeagueRules, error) { return h.Rules.GetLeagueRules(ctx, leagueID, g.RulesVersion) }, "league rules not found")
	if resp != nil {
		return *resp, nil
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

	return responses.JsonResponse(http.StatusCreated, mappers.AtBatToResponse(atbat)), nil
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
		resp = append(resp, mappers.AtBatToResponse(a))
	}

	return responses.JsonResponse(http.StatusOK, resp), nil
}
