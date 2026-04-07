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

type TeamHandlers struct {
	Teams db.TeamRepository
}

func NewTeamHandlers(teams db.TeamRepository) *TeamHandlers {
	return &TeamHandlers{Teams: teams}
}

func (h *TeamHandlers) CreateTeam(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body dto.CreateTeamRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}), nil
	}
	leagueID := req.PathParameters["leagueId"]
	if leagueID == "" || body.Name == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "leagueId and name are required"}), nil
	}
	team := mappers.CreateTeamRequestToDomain(body, uuid.NewString(), leagueID)
	if err := h.Teams.PutTeam(ctx, team); err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}
	return responses.JsonResponse(http.StatusCreated, mappers.TeamToResponse(team)), nil
}

func (h *TeamHandlers) GetTeams(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	leagueID := req.QueryStringParameters["leagueId"]
	if leagueID == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "leagueId query parameter is required"}), nil
	}
	teams, err := h.Teams.ListTeamsByLeague(ctx, leagueID)
	if err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}
	var out []dto.TeamResponse
	for _, t := range teams {
		out = append(out, mappers.TeamToResponse(t))
	}
	return responses.JsonResponse(http.StatusOK, out), nil
}
