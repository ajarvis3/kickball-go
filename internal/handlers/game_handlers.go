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
	"github.com/ajarvis3/kickball-go/internal/mappers"
	"github.com/ajarvis3/kickball-go/pkg/responses"
)

type GameHandlers struct {
	Games db.GameRepository
	Rules db.LeagueRulesRepository
}

func NewGameHandlers(games db.GameRepository, rules db.LeagueRulesRepository) *GameHandlers {
	return &GameHandlers{Games: games, Rules: rules}
}

func (h *GameHandlers) CreateGame(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body dto.CreateGameRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}), nil
	}
	leagueID := req.PathParameters["leagueId"]
	if leagueID == "" || body.HomeTeamID == "" || body.AwayTeamID == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "leagueId, homeTeamId and awayTeamId are required"}), nil
	}
	game := mappers.CreateGameRequestToDomain(body, uuid.NewString(), leagueID)

	// Load league rules to size inning runs and initialize state
	lr, resp := fetchResource(func() (*domain.LeagueRules, error) { return h.Rules.GetLeagueRules(ctx, leagueID, game.RulesVersion) }, "league rules not found for rulesVersion")
	if resp != nil {
		// Treat missing rules as bad request (keeps previous behavior)
		if resp.StatusCode == http.StatusNotFound {
			return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "league rules not found for rulesVersion"}), nil
		}
		return *resp, nil
	}

	game.State.Inning = 1
	game.State.Half = "top"
	game.State.Outs = 0
	if lr.MaxInnings > 0 {
		game.State.InningRuns = make([]int, 2*lr.MaxInnings)
	} else {
		game.State.InningRuns = []int{}
	}

	if err := h.Games.PutGame(ctx, game); err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}
	return responses.JsonResponse(http.StatusCreated, mappers.GameToResponse(game)), nil
}

func (h *GameHandlers) GetGame(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	gameID := req.QueryStringParameters["gameId"]
	leagueID := req.QueryStringParameters["leagueId"]
	if gameID == "" && leagueID == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "gameId or leagueId query parameter is required"}), nil
	}

	if gameID != "" {
		g, resp := fetchResource(func() (*domain.Game, error) { return h.Games.GetGame(ctx, gameID) }, "game not found")
		if resp != nil {
			return *resp, nil
		}
		return responses.JsonResponse(http.StatusOK, mappers.GameToResponse(*g)), nil
	}

	games, err := h.Games.ListGamesByLeague(ctx, leagueID)
	if err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}
	var out []dto.GameResponse
	for _, g := range games {
		out = append(out, mappers.GameToResponse(g))
	}
	return responses.JsonResponse(http.StatusOK, out), nil
}
