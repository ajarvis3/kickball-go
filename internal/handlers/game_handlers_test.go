package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/ajarvis3/kickball-go/internal/domain"
)

type mockGameRepo struct {
	putGameFn           func(ctx context.Context, game domain.Game) error
	getGameFn           func(ctx context.Context, gameID string) (*domain.Game, error)
	listGamesByLeagueFn func(ctx context.Context, leagueID string) ([]domain.Game, error)
}

func (m *mockGameRepo) PutGame(ctx context.Context, game domain.Game) error {
	return m.putGameFn(ctx, game)
}

func (m *mockGameRepo) GetGame(ctx context.Context, gameID string) (*domain.Game, error) {
	return m.getGameFn(ctx, gameID)
}

func (m *mockGameRepo) ListGamesByLeague(ctx context.Context, leagueID string) ([]domain.Game, error) {
	return m.listGamesByLeagueFn(ctx, leagueID)
}

type mockLeagueRulesRepo struct {
	putLeagueRulesFn       func(ctx context.Context, rules domain.LeagueRules) error
	getLeagueRulesFn       func(ctx context.Context, leagueID string, rulesVersion int) (*domain.LeagueRules, error)
	getLatestLeagueRulesFn func(ctx context.Context, leagueID string) (*domain.LeagueRules, error)
}

func (m *mockLeagueRulesRepo) PutLeagueRules(ctx context.Context, rules domain.LeagueRules) error {
	return m.putLeagueRulesFn(ctx, rules)
}

func (m *mockLeagueRulesRepo) GetLeagueRules(ctx context.Context, leagueID string, rulesVersion int) (*domain.LeagueRules, error) {
	return m.getLeagueRulesFn(ctx, leagueID, rulesVersion)
}

func (m *mockLeagueRulesRepo) GetLatestLeagueRules(ctx context.Context, leagueID string) (*domain.LeagueRules, error) {
	if m.getLatestLeagueRulesFn != nil {
		return m.getLatestLeagueRulesFn(ctx, leagueID)
	}
	if m.getLeagueRulesFn != nil {
		return m.getLeagueRulesFn(ctx, leagueID, 0)
	}
	return nil, nil
}

func defaultRules() *domain.LeagueRules {
	return &domain.LeagueRules{
		LeagueID:     "l1",
		RulesVersion: 1,
		MaxStrikes:   3,
		MaxBalls:     4,
		MaxInnings:   7,
	}
}

func TestCreateGameSuccess(t *testing.T) {
	gameRepo := &mockGameRepo{
		putGameFn: func(_ context.Context, _ domain.Game) error { return nil },
	}
	rulesRepo := &mockLeagueRulesRepo{
		getLeagueRulesFn: func(_ context.Context, _ string, _ int) (*domain.LeagueRules, error) {
			return defaultRules(), nil
		},
	}
	h := NewGameHandlers(gameRepo, rulesRepo)
	body, _ := json.Marshal(map[string]string{"homeTeamId": "ht1", "awayTeamId": "at1", "leagueId": "l1"})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.CreateGame(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status = %d; want 201", resp.StatusCode)
	}
}

func TestCreateGameMissingLeagueId(t *testing.T) {
	gameRepo := &mockGameRepo{}
	rulesRepo := &mockLeagueRulesRepo{}
	h := NewGameHandlers(gameRepo, rulesRepo)
	body, _ := json.Marshal(map[string]string{"homeTeamId": "ht1", "awayTeamId": "at1"})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.CreateGame(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestCreateGameRulesNotFound(t *testing.T) {
	gameRepo := &mockGameRepo{}
	rulesRepo := &mockLeagueRulesRepo{
		getLeagueRulesFn: func(_ context.Context, _ string, _ int) (*domain.LeagueRules, error) {
			return nil, nil // not found
		},
	}
	h := NewGameHandlers(gameRepo, rulesRepo)
	body, _ := json.Marshal(map[string]string{"homeTeamId": "ht1", "awayTeamId": "at1", "leagueId": "l1"})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.CreateGame(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestCreateGameRepoError(t *testing.T) {
	gameRepo := &mockGameRepo{
		putGameFn: func(_ context.Context, _ domain.Game) error { return errors.New("db error") },
	}
	rulesRepo := &mockLeagueRulesRepo{
		getLeagueRulesFn: func(_ context.Context, _ string, _ int) (*domain.LeagueRules, error) {
			return defaultRules(), nil
		},
	}
	h := NewGameHandlers(gameRepo, rulesRepo)
	body, _ := json.Marshal(map[string]string{"homeTeamId": "ht1", "awayTeamId": "at1", "leagueId": "l1"})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.CreateGame(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d; want 500", resp.StatusCode)
	}
}

func TestGetGameByGameIdSuccess(t *testing.T) {
	game := &domain.Game{GameID: "g1", LeagueID: "l1", State: domain.GameState{InningRuns: []int{}}}
	gameRepo := &mockGameRepo{
		getGameFn: func(_ context.Context, _ string) (*domain.Game, error) { return game, nil },
	}
	rulesRepo := &mockLeagueRulesRepo{}
	h := NewGameHandlers(gameRepo, rulesRepo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"gameId": "g1"},
	}
	resp, err := h.GetGame(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d; want 200", resp.StatusCode)
	}
}

func TestGetGameByLeagueIdSuccess(t *testing.T) {
	gameRepo := &mockGameRepo{
		listGamesByLeagueFn: func(_ context.Context, _ string) ([]domain.Game, error) {
			return []domain.Game{{GameID: "g1"}}, nil
		},
	}
	rulesRepo := &mockLeagueRulesRepo{}
	h := NewGameHandlers(gameRepo, rulesRepo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"leagueId": "l1"},
	}
	resp, err := h.GetGame(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d; want 200", resp.StatusCode)
	}
}

func TestGetGameMissingParams(t *testing.T) {
	gameRepo := &mockGameRepo{}
	rulesRepo := &mockLeagueRulesRepo{}
	h := NewGameHandlers(gameRepo, rulesRepo)
	req := events.APIGatewayProxyRequest{}
	resp, err := h.GetGame(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestGetGameNotFound(t *testing.T) {
	gameRepo := &mockGameRepo{
		getGameFn: func(_ context.Context, _ string) (*domain.Game, error) { return nil, nil },
	}
	rulesRepo := &mockLeagueRulesRepo{}
	h := NewGameHandlers(gameRepo, rulesRepo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"gameId": "g1"},
	}
	resp, err := h.GetGame(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d; want 404", resp.StatusCode)
	}
}
