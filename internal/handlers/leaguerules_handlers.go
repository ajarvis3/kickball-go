package handlers

import (
    "context"
    "encoding/json"
    "net/http"

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
        LeagueID:            leagueID,
        RulesVersion:        version,
        MaxStrikes:          body.MaxStrikes,
        MaxBalls:            body.MaxBalls,
        MaxFouls:            body.MaxFouls,
        MaxInnings:          body.MaxInnings,
        MercyRunsPerInning:  body.MercyRunsPerInning,
        MercyAppliesLastInning: body.MercyAppliesLastInning,
        GameMercyRuns:       body.GameMercyRuns,
        Metadata:            body.Metadata,
    }

    if err := h.Rules.PutLeagueRules(ctx, rules); err != nil {
        return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), err
    }

    resp := dto.LeagueRulesResponse{
        LeagueID:            rules.LeagueID,
        RulesVersion:        rules.RulesVersion,
        MaxStrikes:          rules.MaxStrikes,
        MaxBalls:            rules.MaxBalls,
        MaxFouls:            rules.MaxFouls,
        MaxInnings:          rules.MaxInnings,
        MercyRunsPerInning:  rules.MercyRunsPerInning,
        MercyAppliesLastInning: rules.MercyAppliesLastInning,
        GameMercyRuns:       rules.GameMercyRuns,
        Metadata:            rules.Metadata,
    }

    return responses.JsonResponse(http.StatusCreated, resp), nil
}
