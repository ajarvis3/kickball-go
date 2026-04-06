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
	team := domain.Team{TeamID: uuid.NewString(), LeagueID: leagueID, Name: body.Name}
	if err := h.Teams.PutTeam(ctx, team); err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), err
	}
	resp := dto.TeamResponse{TeamID: team.TeamID, LeagueID: team.LeagueID, Name: team.Name}
	return responses.JsonResponse(http.StatusCreated, resp), nil
}