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

type AtBatHandlers struct {
	AtBats db.AtBatRepository
}

func NewAtBatHandlers(atbats db.AtBatRepository) *AtBatHandlers {
	return &AtBatHandlers{AtBats: atbats}
}

func (h *AtBatHandlers) RecordAtBat(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body dto.RecordAtBatRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}), nil
	}
	gameID := req.PathParameters["gameId"]
	leagueID := req.PathParameters["leagueId"]
	if gameID == "" || body.PlayerID == "" {
		return responses.JsonResponse(http.StatusBadRequest, map[string]string{"error": "gameId and playerId are required"}), nil
	}
	atbat := domain.AtBat{GameID: gameID, LeagueID: leagueID, TeamID: body.TeamID, PlayerID: body.PlayerID, Seq: 1, Inning: 0, Half: "", Strikes: body.Strikes, Balls: body.Balls, Fouls: body.Fouls, Result: body.Result, RBI: body.RBI, Pitches: body.Pitches}
	if err := h.AtBats.PutAtBat(ctx, atbat); err != nil {
		return responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), err
	}
	resp := dto.AtBatResponse{GameID: atbat.GameID, PlayerID: atbat.PlayerID, TeamID: atbat.TeamID, Seq: atbat.Seq, Inning: atbat.Inning, Half: atbat.Half, Strikes: atbat.Strikes, Balls: atbat.Balls, Fouls: atbat.Fouls, Result: atbat.Result, RBI: atbat.RBI, Pitches: atbat.Pitches}
	return responses.JsonResponse(http.StatusCreated, resp), nil
}