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

type LeagueHandlers struct {
	Leagues db.LeagueRepository
}

func NewLeagueHandlers(leagues db.LeagueRepository) *LeagueHandlers {
	return &LeagueHandlers{Leagues: leagues}
}

func (h *LeagueHandlers) CreateLeague(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO: parse req, call repo, return response
	var body dto.CreateLeagueRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}), nil
	}
	if body.Name == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "league name is required"}), nil
	}
	league := domain.League{
		LeagueID:            uuid.NewString(),
		Name:                body.Name,
		CurrentRulesVersion: 1,
	}
	err := h.Leagues.PutLeague(ctx, league)
	if err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}
	resp := dto.LeagueResponse{
		LeagueID:            league.LeagueID,
		Name:                league.Name,
		CurrentRulesVersion: league.CurrentRulesVersion,
	}
	return responses.JsonResponse(http.StatusCreated, resp), nil
}

func (h *LeagueHandlers) GetLeague(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	leagueID := req.QueryStringParameters["leagueId"]
	if leagueID == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "leagueId query parameter is required"}), nil
	}
	lg, err := h.Leagues.GetLeague(ctx, leagueID)
	if err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}
	if lg == nil {
		return responses.JsonResponse(http.StatusNotFound, map[string]string{"error": "league not found"}), nil
	}
	resp := dto.LeagueResponse{LeagueID: lg.LeagueID, Name: lg.Name, CurrentRulesVersion: lg.CurrentRulesVersion}
	return responses.JsonResponse(http.StatusOK, resp), nil
}
