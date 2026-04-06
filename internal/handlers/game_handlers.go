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

type GameHandlers struct {
	Games db.GameRepository
}

func NewGameHandlers(games db.GameRepository) *GameHandlers {
	return &GameHandlers{Games: games}
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
	game := domain.Game{GameID: uuid.NewString(), LeagueID: leagueID, RulesVersion: 1, HomeTeamID: body.HomeTeamID, AwayTeamID: body.AwayTeamID, State: domain.GameState{}}
	if err := h.Games.PutGame(ctx, game); err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), err
	}
	resp := dto.GameResponse{GameID: game.GameID, LeagueID: game.LeagueID, HomeTeamID: game.HomeTeamID, AwayTeamID: game.AwayTeamID, RulesVersion: game.RulesVersion, State: dto.GameStateDTO{Inning: game.State.Inning, Half: game.State.Half, Outs: game.State.Outs, HomeScore: game.State.HomeScore, AwayScore: game.State.AwayScore}}
	return responses.JsonResponse(http.StatusCreated, resp), nil
}