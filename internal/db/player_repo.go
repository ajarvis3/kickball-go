package db

import (
	"context"

	"github.com/ajarvis3/kickball-go/internal/domain"
)

type PlayerRepository interface {
	PutPlayer(ctx context.Context, player domain.Player) error
	ListPlayersByTeam(ctx context.Context, teamID string) ([]domain.Player, error)
}

type playerRepo struct {
	client *Client
}

func NewPlayerRepository(client *Client) PlayerRepository {
	return &playerRepo{client: client}
}

func (r *playerRepo) PutPlayer(ctx context.Context, player domain.Player) error {
	// TODO
	return nil
}

func (r *playerRepo) ListPlayersByTeam(ctx context.Context, teamID string) ([]domain.Player, error) {
	// TODO
	return nil, nil
}