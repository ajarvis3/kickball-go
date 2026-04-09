package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"

	"github.com/ajarvis3/kickball-go/internal/db"
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/handlers/dto"
	"github.com/ajarvis3/kickball-go/internal/mappers"
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
	league := mappers.CreateLeagueRequestToDomain(body, uuid.NewString())
	err := h.Leagues.PutLeague(ctx, league)
	if err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}
	return responses.JsonResponse(http.StatusCreated, mappers.LeagueToResponse(league)), nil
}

func (h *LeagueHandlers) GetLeague(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	leagueID := req.QueryStringParameters["leagueId"]
	leagueName := req.QueryStringParameters["leagueName"]
	if leagueID == "" {
		// No leagueId provided -> return all leagues (optionally filter by leagueName)
		leagues, err := h.Leagues.ListLeagues(ctx)
		if err != nil {
			return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
		}
		var out []dto.LeagueResponse
		// TODO update this to use a GSI
		if leagueName == "" {
			for _, l := range leagues {
				out = append(out, mappers.LeagueToResponse(l))
			}
			return responses.JsonResponse(http.StatusOK, out), nil
		}

		// filter by name (case-insensitive, substring match)
		q := strings.ToLower(leagueName)
		for _, l := range leagues {
			if strings.Contains(strings.ToLower(l.Name), q) {
				out = append(out, mappers.LeagueToResponse(l))
			}
		}
		return responses.JsonResponse(http.StatusOK, out), nil
	}

	lg, resp := fetchResource(func() (*domain.League, error) { return h.Leagues.GetLeague(ctx, leagueID) }, "league not found")
	if resp != nil {
		return *resp, nil
	}
	return responses.JsonResponse(http.StatusOK, mappers.LeagueToResponse(*lg)), nil
}
