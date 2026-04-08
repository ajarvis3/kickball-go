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

type mockPlayerRepo struct {
	putPlayerFn          func(ctx context.Context, player domain.Player) error
	listPlayersByTeamFn  func(ctx context.Context, teamID string) ([]domain.Player, error)
}

func (m *mockPlayerRepo) PutPlayer(ctx context.Context, player domain.Player) error {
	return m.putPlayerFn(ctx, player)
}

func (m *mockPlayerRepo) ListPlayersByTeam(ctx context.Context, teamID string) ([]domain.Player, error) {
	return m.listPlayersByTeamFn(ctx, teamID)
}

func TestCreatePlayer_Success(t *testing.T) {
	repo := &mockPlayerRepo{
		putPlayerFn: func(_ context.Context, _ domain.Player) error { return nil },
	}
	h := NewPlayerHandlers(repo)
	body, _ := json.Marshal(map[string]interface{}{"name": "Alice", "number": 7, "position": "ss"})
	req := events.APIGatewayProxyRequest{
		Body: string(body),
		PathParameters: map[string]string{
			"teamId":   "t1",
			"leagueId": "l1",
		},
	}
	resp, err := h.CreatePlayer(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status = %d; want 201", resp.StatusCode)
	}
}

func TestCreatePlayer_MissingTeamId(t *testing.T) {
	repo := &mockPlayerRepo{}
	h := NewPlayerHandlers(repo)
	body, _ := json.Marshal(map[string]string{"name": "Alice"})
	req := events.APIGatewayProxyRequest{
		Body:           string(body),
		PathParameters: map[string]string{"leagueId": "l1"},
	}
	resp, err := h.CreatePlayer(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestCreatePlayer_MissingName(t *testing.T) {
	repo := &mockPlayerRepo{}
	h := NewPlayerHandlers(repo)
	body, _ := json.Marshal(map[string]string{"name": ""})
	req := events.APIGatewayProxyRequest{
		Body: string(body),
		PathParameters: map[string]string{
			"teamId":   "t1",
			"leagueId": "l1",
		},
	}
	resp, err := h.CreatePlayer(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestCreatePlayer_RepoError(t *testing.T) {
	repo := &mockPlayerRepo{
		putPlayerFn: func(_ context.Context, _ domain.Player) error { return errors.New("db error") },
	}
	h := NewPlayerHandlers(repo)
	body, _ := json.Marshal(map[string]string{"name": "Alice"})
	req := events.APIGatewayProxyRequest{
		Body: string(body),
		PathParameters: map[string]string{
			"teamId":   "t1",
			"leagueId": "l1",
		},
	}
	resp, err := h.CreatePlayer(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d; want 500", resp.StatusCode)
	}
}

func TestGetPlayers_Success(t *testing.T) {
	repo := &mockPlayerRepo{
		listPlayersByTeamFn: func(_ context.Context, _ string) ([]domain.Player, error) {
			return []domain.Player{{PlayerID: "p1", Name: "Alice"}}, nil
		},
	}
	h := NewPlayerHandlers(repo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"teamId": "t1"},
	}
	resp, err := h.GetPlayers(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d; want 200", resp.StatusCode)
	}
}

func TestGetPlayers_MissingTeamId(t *testing.T) {
	repo := &mockPlayerRepo{}
	h := NewPlayerHandlers(repo)
	req := events.APIGatewayProxyRequest{}
	resp, err := h.GetPlayers(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestGetPlayers_RepoError(t *testing.T) {
	repo := &mockPlayerRepo{
		listPlayersByTeamFn: func(_ context.Context, _ string) ([]domain.Player, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewPlayerHandlers(repo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"teamId": "t1"},
	}
	resp, err := h.GetPlayers(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d; want 500", resp.StatusCode)
	}
}
