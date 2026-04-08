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

type mockTeamRepo struct {
	putTeamFn           func(ctx context.Context, team domain.Team) error
	listTeamsByLeagueFn func(ctx context.Context, leagueID string) ([]domain.Team, error)
}

func (m *mockTeamRepo) PutTeam(ctx context.Context, team domain.Team) error {
	return m.putTeamFn(ctx, team)
}

func (m *mockTeamRepo) ListTeamsByLeague(ctx context.Context, leagueID string) ([]domain.Team, error) {
	return m.listTeamsByLeagueFn(ctx, leagueID)
}

func TestCreateTeam_Success(t *testing.T) {
	repo := &mockTeamRepo{
		putTeamFn: func(_ context.Context, _ domain.Team) error { return nil },
	}
	h := NewTeamHandlers(repo)
	body, _ := json.Marshal(map[string]string{"name": "Tigers"})
	req := events.APIGatewayProxyRequest{
		Body:           string(body),
		PathParameters: map[string]string{"leagueId": "l1"},
	}
	resp, err := h.CreateTeam(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status = %d; want 201", resp.StatusCode)
	}
}

func TestCreateTeam_MissingLeagueId(t *testing.T) {
	repo := &mockTeamRepo{}
	h := NewTeamHandlers(repo)
	body, _ := json.Marshal(map[string]string{"name": "Tigers"})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.CreateTeam(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestCreateTeam_MissingName(t *testing.T) {
	repo := &mockTeamRepo{}
	h := NewTeamHandlers(repo)
	body, _ := json.Marshal(map[string]string{"name": ""})
	req := events.APIGatewayProxyRequest{
		Body:           string(body),
		PathParameters: map[string]string{"leagueId": "l1"},
	}
	resp, err := h.CreateTeam(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestCreateTeam_RepoError(t *testing.T) {
	repo := &mockTeamRepo{
		putTeamFn: func(_ context.Context, _ domain.Team) error { return errors.New("db error") },
	}
	h := NewTeamHandlers(repo)
	body, _ := json.Marshal(map[string]string{"name": "Tigers"})
	req := events.APIGatewayProxyRequest{
		Body:           string(body),
		PathParameters: map[string]string{"leagueId": "l1"},
	}
	resp, err := h.CreateTeam(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d; want 500", resp.StatusCode)
	}
}

func TestGetTeams_Success(t *testing.T) {
	repo := &mockTeamRepo{
		listTeamsByLeagueFn: func(_ context.Context, _ string) ([]domain.Team, error) {
			return []domain.Team{{TeamID: "t1", Name: "Tigers"}}, nil
		},
	}
	h := NewTeamHandlers(repo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"leagueId": "l1"},
	}
	resp, err := h.GetTeams(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d; want 200", resp.StatusCode)
	}
}

func TestGetTeams_MissingLeagueId(t *testing.T) {
	repo := &mockTeamRepo{}
	h := NewTeamHandlers(repo)
	req := events.APIGatewayProxyRequest{}
	resp, err := h.GetTeams(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestGetTeams_RepoError(t *testing.T) {
	repo := &mockTeamRepo{
		listTeamsByLeagueFn: func(_ context.Context, _ string) ([]domain.Team, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewTeamHandlers(repo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"leagueId": "l1"},
	}
	resp, err := h.GetTeams(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d; want 500", resp.StatusCode)
	}
}
