package db

import (
	"context"

	"github.com/ajarvis3/kickball-go/internal/domain"
)

type GameRepository interface {
	PutGame(ctx context.Context, game domain.Game) error
	GetGame(ctx context.Context, gameID string) (*domain.Game, error)
	ListGamesByLeague(ctx context.Context, leagueID string) ([]domain.Game, error)
}

type gameRepo struct {
	client *Client
}

func NewGameRepository(client *Client) GameRepository {
	return &gameRepo{client: client}
}

func (r *gameRepo) PutGame(ctx context.Context, game domain.Game) error {
	// TODO
	return nil
}

func (r *gameRepo) GetGame(ctx context.Context, gameID string) (*domain.Game, error) {
	// TODO
	return nil, nil
}

func (r *gameRepo) ListGamesByLeague(ctx context.Context, leagueID string) ([]domain.Game, error) {
	// TODO
	return nil, nil
}