package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

	"github.com/ajarvis3/kickball-go/internal/data/domain"
	dto "github.com/ajarvis3/kickball-go/internal/data/dto"
	"github.com/ajarvis3/kickball-go/internal/db"
	"github.com/ajarvis3/kickball-go/internal/mappers"
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

	leagueID := body.LeagueID
	if leagueID == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "leagueId is required"}), nil
	}

	// default rules version if not provided
	version := body.RulesVersion
	if version == 0 {
		version = 1
	}

	rules := mappers.CreateLeagueRulesRequestToDomain(body, leagueID, version)

	if err := h.Rules.PutLeagueRules(ctx, rules); err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}

	return responses.JsonResponse(http.StatusCreated, mappers.LeagueRulesToResponse(rules)), nil
}

func (h *LeagueRulesHandlers) GetLeagueRules(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Expect query params: leagueId and version
	leagueID := req.QueryStringParameters["leagueId"]
	versionStr := req.QueryStringParameters["version"]
	if leagueID == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "leagueId query parameter is required"}), nil
	}

	// If version not provided, return the latest rules for the league
	if versionStr == "" {
		rules, resp := fetchResource(func() (*domain.LeagueRules, error) { return h.Rules.GetLatestLeagueRules(ctx, leagueID) }, "league rules not found")
		if resp != nil {
			return *resp, nil
		}
		return responses.JsonResponse(http.StatusOK, mappers.LeagueRulesToResponse(*rules)), nil
	}

	v, err := strconv.Atoi(versionStr)
	if err != nil {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "invalid version"}), nil
	}

	rules, resp := fetchResource(func() (*domain.LeagueRules, error) { return h.Rules.GetLeagueRules(ctx, leagueID, v) }, "league rules not found")
	if resp != nil {
		return *resp, nil
	}

	return responses.JsonResponse(http.StatusOK, mappers.LeagueRulesToResponse(*rules)), nil
}
