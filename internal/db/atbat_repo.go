package db

import (
	"context"

	"github.com/ajarvis3/kickball-go/internal/domain"
)

type AtBatRepository interface {
	PutAtBat(ctx context.Context, atbat domain.AtBat) error
	ListAtBatsByGame(ctx context.Context, gameID string) ([]domain.AtBat, error)
	ListAtBatsByPlayer(ctx context.Context, playerID string) ([]domain.AtBat, error)
}

type atBatRepo struct {
	client *Client
}

func NewAtBatRepository(client *Client) AtBatRepository {
	return &atBatRepo{client: client}
}

func (r *atBatRepo) PutAtBat(ctx context.Context, atbat domain.AtBat) error {
	// TODO
	return nil
}

func (r *atBatRepo) ListAtBatsByGame(ctx context.Context, gameID string) ([]domain.AtBat, error) {
	// TODO
	return nil, nil
}

func (r *atBatRepo) ListAtBatsByPlayer(ctx context.Context, playerID string) ([]domain.AtBat, error) {
	// TODO
	return nil, nil
}