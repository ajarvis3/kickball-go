package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

	"github.com/ajarvis3/kickball-go/internal/db"
	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/handlers/dto"
	"github.com/ajarvis3/kickball-go/pkg/responses"
)

type LeagueRulesHandlers struct {
	Rules db.LeagueRulesRepository
}

func NewLeagueRulesHandlers(r db.LeagueRulesRepository) *LeagueRulesHandlers {
	return &LeagueRulesHandlers{Rules: r}
}

func (h *LeagueRulesHandlers) CreateLeagueRules(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body dto.CreateLeagueRulesRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}), nil
	}

	leagueID := req.PathParameters["leagueId"]
	if leagueID == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "leagueId is required"}), nil
	}

	// default rules version if not provided
	version := body.RulesVersion
	if version == 0 {
		version = 1
	}

	rules := domain.LeagueRules{
		LeagueID:               leagueID,
		RulesVersion:           version,
		MaxStrikes:             body.MaxStrikes,
		MaxBalls:               body.MaxBalls,
		MaxFouls:               body.MaxFouls,
		MaxInnings:             body.MaxInnings,
		MercyRunsPerInning:     body.MercyRunsPerInning,
		MercyAppliesLastInning: body.MercyAppliesLastInning,
		GameMercyRuns:          body.GameMercyRuns,
		Metadata:               body.Metadata,
	}

	if err := h.Rules.PutLeagueRules(ctx, rules); err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}

	resp := dto.LeagueRulesResponse{
		LeagueID:               rules.LeagueID,
		RulesVersion:           rules.RulesVersion,
		MaxStrikes:             rules.MaxStrikes,
		MaxBalls:               rules.MaxBalls,
		MaxFouls:               rules.MaxFouls,
		MaxInnings:             rules.MaxInnings,
		MercyRunsPerInning:     rules.MercyRunsPerInning,
		MercyAppliesLastInning: rules.MercyAppliesLastInning,
		GameMercyRuns:          rules.GameMercyRuns,
		Metadata:               rules.Metadata,
	}

	return responses.JsonResponse(http.StatusCreated, resp), nil
}

func (h *LeagueRulesHandlers) GetLeagueRules(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Expect query params: leagueId and version
	leagueID := req.QueryStringParameters["leagueId"]
	versionStr := req.QueryStringParameters["version"]
	if leagueID == "" || versionStr == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "leagueId and version query parameters are required"}), nil
	}

	v, err := strconv.Atoi(versionStr)
	if err != nil {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "invalid version"}), nil
	}

	rules, resp := fetchResource(func() (*domain.LeagueRules, error) { return h.Rules.GetLeagueRules(ctx, leagueID, v) }, "league rules not found")
	if resp != nil {
		return *resp, nil
	}

	// map domain -> dto
	respDto := dto.LeagueRulesResponse{
		LeagueID:               rules.LeagueID,
		RulesVersion:           rules.RulesVersion,
		MaxStrikes:             rules.MaxStrikes,
		MaxBalls:               rules.MaxBalls,
		MaxFouls:               rules.MaxFouls,
		MaxInnings:             rules.MaxInnings,
		MercyRunsPerInning:     rules.MercyRunsPerInning,
		MercyAppliesLastInning: rules.MercyAppliesLastInning,
		GameMercyRuns:          rules.GameMercyRuns,
		Metadata:               rules.Metadata,
	}

	return responses.JsonResponse(http.StatusOK, respDto), nil
}
