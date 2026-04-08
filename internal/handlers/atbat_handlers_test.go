package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/services"
)

type mockAtBatRepo struct {
	putAtBatFn              func(ctx context.Context, atbat domain.AtBat) error
	putAtBatAndUpdateGameFn func(ctx context.Context, atbat domain.AtBat, updatedGame domain.Game) error
	listAtBatsByGameFn      func(ctx context.Context, gameID string) ([]domain.AtBat, error)
	listAtBatsByPlayerFn    func(ctx context.Context, playerID string) ([]domain.AtBat, error)
}

func (m *mockAtBatRepo) PutAtBat(ctx context.Context, atbat domain.AtBat) error {
	return m.putAtBatFn(ctx, atbat)
}

func (m *mockAtBatRepo) PutAtBatAndUpdateGame(ctx context.Context, atbat domain.AtBat, updatedGame domain.Game) error {
	return m.putAtBatAndUpdateGameFn(ctx, atbat, updatedGame)
}

func (m *mockAtBatRepo) ListAtBatsByGame(ctx context.Context, gameID string) ([]domain.AtBat, error) {
	return m.listAtBatsByGameFn(ctx, gameID)
}

func (m *mockAtBatRepo) ListAtBatsByPlayer(ctx context.Context, playerID string) ([]domain.AtBat, error) {
	return m.listAtBatsByPlayerFn(ctx, playerID)
}

// mockGameRepoForAtBat is a separate type from mockGameRepo to avoid duplicate in package.
type mockGameRepoForAtBat struct {
	putGameFn           func(ctx context.Context, game domain.Game) error
	getGameFn           func(ctx context.Context, gameID string) (*domain.Game, error)
	listGamesByLeagueFn func(ctx context.Context, leagueID string) ([]domain.Game, error)
}

func (m *mockGameRepoForAtBat) PutGame(ctx context.Context, game domain.Game) error {
	return m.putGameFn(ctx, game)
}

func (m *mockGameRepoForAtBat) GetGame(ctx context.Context, gameID string) (*domain.Game, error) {
	return m.getGameFn(ctx, gameID)
}

func (m *mockGameRepoForAtBat) ListGamesByLeague(ctx context.Context, leagueID string) ([]domain.Game, error) {
	return m.listGamesByLeagueFn(ctx, leagueID)
}

// mockLeagueRulesRepoForAtBat is a separate type from mockLeagueRulesRepo/LR to avoid duplicates.
type mockLeagueRulesRepoForAtBat struct {
	putLeagueRulesFn       func(ctx context.Context, rules domain.LeagueRules) error
	getLeagueRulesFn       func(ctx context.Context, leagueID string, rulesVersion int) (*domain.LeagueRules, error)
	getLatestLeagueRulesFn func(ctx context.Context, leagueID string) (*domain.LeagueRules, error)
}

func (m *mockLeagueRulesRepoForAtBat) PutLeagueRules(ctx context.Context, rules domain.LeagueRules) error {
	return m.putLeagueRulesFn(ctx, rules)
}

func (m *mockLeagueRulesRepoForAtBat) GetLeagueRules(ctx context.Context, leagueID string, rulesVersion int) (*domain.LeagueRules, error) {
	return m.getLeagueRulesFn(ctx, leagueID, rulesVersion)
}

func (m *mockLeagueRulesRepoForAtBat) GetLatestLeagueRules(ctx context.Context, leagueID string) (*domain.LeagueRules, error) {
	if m.getLatestLeagueRulesFn != nil {
		return m.getLatestLeagueRulesFn(ctx, leagueID)
	}
	// Fallback to GetLeagueRules with version 0 if not provided in test
	if m.getLeagueRulesFn != nil {
		return m.getLeagueRulesFn(ctx, leagueID, 0)
	}
	return nil, nil
}

func newTestGame() *domain.Game {
	return &domain.Game{
		GameID:       "g1",
		LeagueID:     "l1",
		HomeTeamID:   "ht1",
		AwayTeamID:   "at1",
		RulesVersion: 1,
		State: domain.GameState{
			Inning:     1,
			Half:       "top",
			Outs:       0,
			InningRuns: make([]int, 14), // 2 * 7
		},
	}
}

func newTestRules() *domain.LeagueRules {
	return &domain.LeagueRules{
		LeagueID:     "l1",
		RulesVersion: 1,
		MaxStrikes:   3,
		MaxBalls:     4,
		MaxInnings:   7,
	}
}

func TestRecordAtBatSuccess(t *testing.T) {
	atbatRepo := &mockAtBatRepo{
		putAtBatAndUpdateGameFn: func(_ context.Context, _ domain.AtBat, _ domain.Game) error { return nil },
	}
	gameRepo := &mockGameRepoForAtBat{
		getGameFn: func(_ context.Context, _ string) (*domain.Game, error) { return newTestGame(), nil },
	}
	rulesRepo := &mockLeagueRulesRepoForAtBat{
		getLeagueRulesFn: func(_ context.Context, _ string, _ int) (*domain.LeagueRules, error) {
			return newTestRules(), nil
		},
	}
	engine := services.NewGameEngine(services.NewRulesEngine())
	h := NewAtBatHandlers(atbatRepo, gameRepo, rulesRepo, engine)

	body, _ := json.Marshal(map[string]interface{}{
		"gameId":   "g1",
		"leagueId": "l1",
		"playerId": "p1",
		"teamId":   "ht1",
		"result":   "out",
		"strikes":  1,
		"balls":    2,
	})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.RecordAtBat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status = %d; want 201", resp.StatusCode)
	}
}

func TestRecordAtBatMissingGameId(t *testing.T) {
	atbatRepo := &mockAtBatRepo{}
	gameRepo := &mockGameRepoForAtBat{}
	rulesRepo := &mockLeagueRulesRepoForAtBat{}
	engine := services.NewGameEngine(services.NewRulesEngine())
	h := NewAtBatHandlers(atbatRepo, gameRepo, rulesRepo, engine)

	body, _ := json.Marshal(map[string]interface{}{"playerId": "p1", "result": "out"})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.RecordAtBat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestRecordAtBatMissingPlayerId(t *testing.T) {
	atbatRepo := &mockAtBatRepo{}
	gameRepo := &mockGameRepoForAtBat{}
	rulesRepo := &mockLeagueRulesRepoForAtBat{}
	engine := services.NewGameEngine(services.NewRulesEngine())
	h := NewAtBatHandlers(atbatRepo, gameRepo, rulesRepo, engine)

	body, _ := json.Marshal(map[string]interface{}{"result": "out"})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.RecordAtBat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestRecordAtBatGameNotFound(t *testing.T) {
	atbatRepo := &mockAtBatRepo{}
	gameRepo := &mockGameRepoForAtBat{
		getGameFn: func(_ context.Context, _ string) (*domain.Game, error) { return nil, nil },
	}
	rulesRepo := &mockLeagueRulesRepoForAtBat{}
	engine := services.NewGameEngine(services.NewRulesEngine())
	h := NewAtBatHandlers(atbatRepo, gameRepo, rulesRepo, engine)

	body, _ := json.Marshal(map[string]interface{}{"gameId": "g1", "leagueId": "l1", "playerId": "p1", "result": "out"})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.RecordAtBat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d; want 404", resp.StatusCode)
	}
}

func TestRecordAtBatPersistError(t *testing.T) {
	atbatRepo := &mockAtBatRepo{
		putAtBatAndUpdateGameFn: func(_ context.Context, _ domain.AtBat, _ domain.Game) error {
			return errors.New("db error")
		},
	}
	gameRepo := &mockGameRepoForAtBat{
		getGameFn: func(_ context.Context, _ string) (*domain.Game, error) { return newTestGame(), nil },
	}
	rulesRepo := &mockLeagueRulesRepoForAtBat{
		getLeagueRulesFn: func(_ context.Context, _ string, _ int) (*domain.LeagueRules, error) {
			return newTestRules(), nil
		},
	}
	engine := services.NewGameEngine(services.NewRulesEngine())
	h := NewAtBatHandlers(atbatRepo, gameRepo, rulesRepo, engine)

	body, _ := json.Marshal(map[string]interface{}{"gameId": "g1", "leagueId": "l1", "playerId": "p1", "teamId": "ht1", "result": "out"})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.RecordAtBat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d; want 500", resp.StatusCode)
	}
}

func TestGetAtBatsByGameId(t *testing.T) {
	atbatRepo := &mockAtBatRepo{
		listAtBatsByGameFn: func(_ context.Context, _ string) ([]domain.AtBat, error) {
			return []domain.AtBat{{GameID: "g1", PlayerID: "p1"}}, nil
		},
	}
	gameRepo := &mockGameRepoForAtBat{}
	rulesRepo := &mockLeagueRulesRepoForAtBat{}
	engine := services.NewGameEngine(services.NewRulesEngine())
	h := NewAtBatHandlers(atbatRepo, gameRepo, rulesRepo, engine)

	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"gameId": "g1"},
	}
	resp, err := h.GetAtBats(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d; want 200", resp.StatusCode)
	}
}

func TestGetAtBatsByPlayerId(t *testing.T) {
	atbatRepo := &mockAtBatRepo{
		listAtBatsByPlayerFn: func(_ context.Context, _ string) ([]domain.AtBat, error) {
			return []domain.AtBat{{GameID: "g1", PlayerID: "p1"}}, nil
		},
	}
	gameRepo := &mockGameRepoForAtBat{}
	rulesRepo := &mockLeagueRulesRepoForAtBat{}
	engine := services.NewGameEngine(services.NewRulesEngine())
	h := NewAtBatHandlers(atbatRepo, gameRepo, rulesRepo, engine)

	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"playerId": "p1"},
	}
	resp, err := h.GetAtBats(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d; want 200", resp.StatusCode)
	}
}

func TestGetAtBatsMissingParams(t *testing.T) {
	atbatRepo := &mockAtBatRepo{}
	gameRepo := &mockGameRepoForAtBat{}
	rulesRepo := &mockLeagueRulesRepoForAtBat{}
	engine := services.NewGameEngine(services.NewRulesEngine())
	h := NewAtBatHandlers(atbatRepo, gameRepo, rulesRepo, engine)

	req := events.APIGatewayProxyRequest{}
	resp, err := h.GetAtBats(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestGetAtBatsRepoError(t *testing.T) {
	atbatRepo := &mockAtBatRepo{
		listAtBatsByGameFn: func(_ context.Context, _ string) ([]domain.AtBat, error) {
			return nil, errors.New("db error")
		},
	}
	gameRepo := &mockGameRepoForAtBat{}
	rulesRepo := &mockLeagueRulesRepoForAtBat{}
	engine := services.NewGameEngine(services.NewRulesEngine())
	h := NewAtBatHandlers(atbatRepo, gameRepo, rulesRepo, engine)

	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"gameId": "g1"},
	}
	resp, err := h.GetAtBats(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d; want 500", resp.StatusCode)
	}
}
