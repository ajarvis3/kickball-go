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

type mockLeagueRepo struct {
	putLeagueFn func(ctx context.Context, league domain.League) error
	getLeagueFn func(ctx context.Context, leagueID string) (*domain.League, error)
}

func (m *mockLeagueRepo) PutLeague(ctx context.Context, league domain.League) error {
	return m.putLeagueFn(ctx, league)
}

func (m *mockLeagueRepo) GetLeague(ctx context.Context, leagueID string) (*domain.League, error) {
	return m.getLeagueFn(ctx, leagueID)
}

func TestCreateLeague_Success(t *testing.T) {
	repo := &mockLeagueRepo{
		putLeagueFn: func(_ context.Context, _ domain.League) error { return nil },
	}
	h := NewLeagueHandlers(repo)
	body, _ := json.Marshal(map[string]string{"name": "Test League"})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.CreateLeague(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status = %d; want 201", resp.StatusCode)
	}
}

func TestCreateLeague_BadJSON(t *testing.T) {
	repo := &mockLeagueRepo{}
	h := NewLeagueHandlers(repo)
	req := events.APIGatewayProxyRequest{Body: "not json"}
	resp, err := h.CreateLeague(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestCreateLeague_EmptyName(t *testing.T) {
	repo := &mockLeagueRepo{}
	h := NewLeagueHandlers(repo)
	body, _ := json.Marshal(map[string]string{"name": ""})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.CreateLeague(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestCreateLeague_RepoError(t *testing.T) {
	repo := &mockLeagueRepo{
		putLeagueFn: func(_ context.Context, _ domain.League) error { return errors.New("db error") },
	}
	h := NewLeagueHandlers(repo)
	body, _ := json.Marshal(map[string]string{"name": "League"})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.CreateLeague(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d; want 500", resp.StatusCode)
	}
}

func TestGetLeague_Success(t *testing.T) {
	league := &domain.League{LeagueID: "l1", Name: "My League", CurrentRulesVersion: 1}
	repo := &mockLeagueRepo{
		getLeagueFn: func(_ context.Context, _ string) (*domain.League, error) { return league, nil },
	}
	h := NewLeagueHandlers(repo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"leagueId": "l1"},
	}
	resp, err := h.GetLeague(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d; want 200", resp.StatusCode)
	}
}

func TestGetLeague_MissingLeagueId(t *testing.T) {
	repo := &mockLeagueRepo{}
	h := NewLeagueHandlers(repo)
	req := events.APIGatewayProxyRequest{}
	resp, err := h.GetLeague(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestGetLeague_NotFound(t *testing.T) {
	repo := &mockLeagueRepo{
		getLeagueFn: func(_ context.Context, _ string) (*domain.League, error) { return nil, nil },
	}
	h := NewLeagueHandlers(repo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"leagueId": "l1"},
	}
	resp, err := h.GetLeague(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d; want 404", resp.StatusCode)
	}
}

func TestGetLeague_RepoError(t *testing.T) {
	repo := &mockLeagueRepo{
		getLeagueFn: func(_ context.Context, _ string) (*domain.League, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewLeagueHandlers(repo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"leagueId": "l1"},
	}
	resp, err := h.GetLeague(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d; want 500", resp.StatusCode)
	}
}
