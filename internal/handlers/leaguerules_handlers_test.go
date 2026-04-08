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

// mockLeagueRulesRepoForHandlers is a local mock (distinct name to avoid conflicts within same package).
type mockLeagueRulesRepoForHandlers struct {
	putLeagueRulesFn       func(ctx context.Context, rules domain.LeagueRules) error
	getLeagueRulesFn       func(ctx context.Context, leagueID string, rulesVersion int) (*domain.LeagueRules, error)
	getLatestLeagueRulesFn func(ctx context.Context, leagueID string) (*domain.LeagueRules, error)
}

func (m *mockLeagueRulesRepoForHandlers) PutLeagueRules(ctx context.Context, rules domain.LeagueRules) error {
	return m.putLeagueRulesFn(ctx, rules)
}

func (m *mockLeagueRulesRepoForHandlers) GetLeagueRules(ctx context.Context, leagueID string, rulesVersion int) (*domain.LeagueRules, error) {
	return m.getLeagueRulesFn(ctx, leagueID, rulesVersion)
}

func (m *mockLeagueRulesRepoForHandlers) GetLatestLeagueRules(ctx context.Context, leagueID string) (*domain.LeagueRules, error) {
	if m.getLatestLeagueRulesFn != nil {
		return m.getLatestLeagueRulesFn(ctx, leagueID)
	}
	if m.getLeagueRulesFn != nil {
		return m.getLeagueRulesFn(ctx, leagueID, 0)
	}
	return nil, nil
}

func TestCreateLeagueRulesSuccess(t *testing.T) {
	repo := &mockLeagueRulesRepoForHandlers{
		putLeagueRulesFn: func(_ context.Context, _ domain.LeagueRules) error { return nil },
	}
	h := NewLeagueRulesHandlers(repo)
	body, _ := json.Marshal(map[string]interface{}{
		"maxStrikes": 3,
		"maxBalls":   4,
		"maxInnings": 7,
		"leagueId":   "l1",
	})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.CreateLeagueRules(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status = %d; want 201", resp.StatusCode)
	}
}

func TestCreateLeagueRulesMissingLeagueId(t *testing.T) {
	repo := &mockLeagueRulesRepoForHandlers{}
	h := NewLeagueRulesHandlers(repo)
	body, _ := json.Marshal(map[string]interface{}{"maxStrikes": 3})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.CreateLeagueRules(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestCreateLeagueRulesRepoError(t *testing.T) {
	repo := &mockLeagueRulesRepoForHandlers{
		putLeagueRulesFn: func(_ context.Context, _ domain.LeagueRules) error { return errors.New("db error") },
	}
	h := NewLeagueRulesHandlers(repo)
	body, _ := json.Marshal(map[string]interface{}{"maxStrikes": 3, "leagueId": "l1"})
	req := events.APIGatewayProxyRequest{Body: string(body)}
	resp, err := h.CreateLeagueRules(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d; want 500", resp.StatusCode)
	}
}

func TestGetLeagueRulesSuccess(t *testing.T) {
	rules := &domain.LeagueRules{LeagueID: "l1", RulesVersion: 1, MaxStrikes: 3}
	repo := &mockLeagueRulesRepoForHandlers{
		getLeagueRulesFn: func(_ context.Context, _ string, _ int) (*domain.LeagueRules, error) {
			return rules, nil
		},
	}
	h := NewLeagueRulesHandlers(repo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"leagueId": "l1", "version": "1"},
	}
	resp, err := h.GetLeagueRules(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d; want 200", resp.StatusCode)
	}
}

func TestGetLeagueRulesMissingParams(t *testing.T) {
	repo := &mockLeagueRulesRepoForHandlers{}
	h := NewLeagueRulesHandlers(repo)
	req := events.APIGatewayProxyRequest{}
	resp, err := h.GetLeagueRules(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}

func TestGetLeagueRulesNotFound(t *testing.T) {
	repo := &mockLeagueRulesRepoForHandlers{
		getLeagueRulesFn: func(_ context.Context, _ string, _ int) (*domain.LeagueRules, error) {
			return nil, nil
		},
	}
	h := NewLeagueRulesHandlers(repo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"leagueId": "l1", "version": "1"},
	}
	resp, err := h.GetLeagueRules(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d; want 404", resp.StatusCode)
	}
}

func TestGetLeagueRulesRepoError(t *testing.T) {
	repo := &mockLeagueRulesRepoForHandlers{
		getLeagueRulesFn: func(_ context.Context, _ string, _ int) (*domain.LeagueRules, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewLeagueRulesHandlers(repo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"leagueId": "l1", "version": "1"},
	}
	resp, err := h.GetLeagueRules(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d; want 500", resp.StatusCode)
	}
}

func TestGetLeagueRulesInvalidVersion(t *testing.T) {
	repo := &mockLeagueRulesRepoForHandlers{}
	h := NewLeagueRulesHandlers(repo)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"leagueId": "l1", "version": "abc"},
	}
	resp, err := h.GetLeagueRules(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d; want 400", resp.StatusCode)
	}
}
