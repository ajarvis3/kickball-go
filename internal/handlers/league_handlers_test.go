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
	putLeagueFn   func(ctx context.Context, league domain.League) error
	getLeagueFn   func(ctx context.Context, leagueID string) (*domain.League, error)
	listLeaguesFn func(ctx context.Context) ([]domain.League, error)
	listLeaguesByNameFn func(ctx context.Context, namePrefix string) ([]domain.League, error)
}

func (m *mockLeagueRepo) PutLeague(ctx context.Context, league domain.League) error {
	return m.putLeagueFn(ctx, league)
}

func (m *mockLeagueRepo) GetLeague(ctx context.Context, leagueID string) (*domain.League, error) {
	return m.getLeagueFn(ctx, leagueID)
}

func (m *mockLeagueRepo) ListLeagues(ctx context.Context) ([]domain.League, error) {
	return m.listLeaguesFn(ctx)
}

func (m *mockLeagueRepo) ListLeaguesByName(ctx context.Context, namePrefix string) ([]domain.League, error) {
	if m.listLeaguesByNameFn != nil {
		return m.listLeaguesByNameFn(ctx, namePrefix)
	}
	if m.listLeaguesFn != nil {
		return m.listLeaguesFn(ctx)
	}
	return []domain.League{}, nil
}

func TestCreateLeagueSuccess(t *testing.T) {
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

func TestCreateLeagueBadJSON(t *testing.T) {
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

func TestCreateLeagueEmptyName(t *testing.T) {
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

func TestCreateLeagueRepoError(t *testing.T) {
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

func TestGetLeagueSuccess(t *testing.T) {
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

func TestGetLeagueMissingLeagueId(t *testing.T) {
	repo := &mockLeagueRepo{listLeaguesFn: func(_ context.Context) ([]domain.League, error) { return []domain.League{}, nil }}
	h := NewLeagueHandlers(repo)
	req := events.APIGatewayProxyRequest{}
	resp, err := h.GetLeague(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d; want 200", resp.StatusCode)
	}
}

func TestGetLeagueNotFound(t *testing.T) {
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

func TestGetLeagueRepoError(t *testing.T) {
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
