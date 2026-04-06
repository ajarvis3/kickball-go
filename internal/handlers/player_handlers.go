package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"

	"github.com/ajarvis3/kickball-go/internal/db"
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/handlers/dto"
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
	teamID := req.PathParameters["teamId"]
	leagueID := req.PathParameters["leagueId"]
	if teamID == "" || leagueID == "" || body.Name == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "teamId, leagueId and name are required"}), nil
	}
	player := domain.Player{PlayerID: uuid.NewString(), TeamID: teamID, LeagueID: leagueID, Name: body.Name, Number: body.Number, Position: body.Position}
	if err := h.Players.PutPlayer(ctx, player); err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), err
	}
	resp := dto.PlayerResponse{PlayerID: player.PlayerID, TeamID: player.TeamID, LeagueID: player.LeagueID, Name: player.Name, Number: player.Number, Position: player.Position}
	return responses.JsonResponse(http.StatusCreated, resp), nil
}