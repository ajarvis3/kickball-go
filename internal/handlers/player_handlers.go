package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"

	"github.com/ajarvis3/kickball-go/internal/db"
	"github.com/ajarvis3/kickball-go/internal/handlers/dto"
	"github.com/ajarvis3/kickball-go/internal/mappers"
	"github.com/ajarvis3/kickball-go/pkg/responses"
)

type PlayerHandlers struct {
	Players db.PlayerRepository
}

func NewPlayerHandlers(players db.PlayerRepository) *PlayerHandlers {
	return &PlayerHandlers{Players: players}
}

func (h *PlayerHandlers) CreatePlayer(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body dto.CreatePlayerRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}), nil
	}
	// Accept teamId and leagueId from request body instead of path parameters
	teamID := body.TeamID
	leagueID := body.LeagueID
	if teamID == "" || leagueID == "" || body.Name == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "teamId, leagueId and name are required"}), nil
	}
	player := mappers.CreatePlayerRequestToDomain(body, uuid.NewString(), teamID, leagueID)
	if err := h.Players.PutPlayer(ctx, player); err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}
	return responses.JsonResponse(http.StatusCreated, mappers.PlayerToResponse(player)), nil
}

func (h *PlayerHandlers) GetPlayers(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	teamID := req.QueryStringParameters["teamId"]
	if teamID == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "teamId query parameter is required"}), nil
	}
	players, err := h.Players.ListPlayersByTeam(ctx, teamID)
	if err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}
	var out []dto.PlayerResponse
	for _, p := range players {
		out = append(out, mappers.PlayerToResponse(p))
	}
	return responses.JsonResponse(http.StatusOK, out), nil
}
