package db

import (
	"context"

	"github.com/ajarvis3/kickball-go/internal/domain"
)

type TeamRepository interface {
	PutTeam(ctx context.Context, team domain.Team) error
	ListTeamsByLeague(ctx context.Context, leagueID string) ([]domain.Team, error)
}

type teamRepo struct {
	client *Client
}

func NewTeamRepository(client *Client) TeamRepository {
	return &teamRepo{client: client}
}

func (r *teamRepo) PutTeam(ctx context.Context, team domain.Team) error {
	// TODO
	return nil
}

func (r *teamRepo) ListTeamsByLeague(ctx context.Context, leagueID string) ([]domain.Team, error) {
	// TODO
	return nil, nil
}