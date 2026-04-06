package services

import "github.com/ajarvis3/kickball-go/internal/domain"

type GameEngine struct {
	// TODO: maybe rules, logger, etc.
}

func NewGameEngine() *GameEngine {
	return &GameEngine{}
}

func (e *GameEngine) ApplyAtBat(game domain.Game, atbat domain.AtBat) (domain.Game, error) {
	// TODO: update game state based on at-bat
	return game, nil
}