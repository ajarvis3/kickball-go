package services

import (
	"strings"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/pkg/apperrors"
)

type GameEngine struct {
	Rules *RulesEngine
}

func NewGameEngine(rules *RulesEngine) *GameEngine {
	return &GameEngine{Rules: rules}
}

// ApplyAtBat applies the at-bat to the game state using the provided league rules.
// It updates outs, inning/half transitions, and scores based on the at-bat result and RBI.
func (e *GameEngine) ApplyAtBat(game domain.Game, rules domain.LeagueRules, atbat domain.AtBat) (domain.Game, error) {
	if e == nil || e.Rules == nil {
		return game, apperrors.ErrInternalError
	}

	// Validate counts against rules first
	if err := e.Rules.ValidateAtBat(game, rules, atbat); err != nil {
		return game, err
	}

	// Initialize inning/half if not started
	if game.State.Inning == 0 {
		game.State.Inning = 1
		game.State.Half = "top"
		game.State.Outs = 0
		game.State.InningRuns = []int{}
	}

	res := strings.ToLower(strings.TrimSpace(atbat.Result))

	// Score runs for hitting/walks/etc when RBI provided
	switch res {
	case "single", "double", "triple", "homerun", "walk", "sacrifice", "error", "fielderschoice":
		if atbat.RBI > 0 {
			if atbat.TeamID == game.HomeTeamID {
				game.State.HomeScore += atbat.RBI
			} else {
				game.State.AwayScore += atbat.RBI
			}
			// Ensure inning runs slice is long enough for current inning
			idx := game.State.Inning - 1
			if idx < 0 {
				idx = 0
			}
			if len(game.State.InningRuns) <= idx {
				// extend with zeros
				needed := idx - len(game.State.InningRuns) + 1
				for i := 0; i < needed; i++ {
					game.State.InningRuns = append(game.State.InningRuns, 0)
				}
			}
			game.State.InningRuns[idx] += atbat.RBI
		}
	}

	// Handle outs
	switch res {
	case "out", "strikeout":
		game.State.Outs += 1
	case "doubleplay":
		game.State.Outs += 2
	case "tripleplay":
		game.State.Outs += 3
	}

	// Advance inning/half when outs reach 3 or more
	for game.State.Outs >= 3 {
		game.State.Outs -= 3
		if game.State.Half == "top" {
			game.State.Half = "bottom"
		} else {
			game.State.Half = "top"
			game.State.Inning += 1
		}
	}

	// Basic mercy rule check (game-level)
	if rules.GameMercyRuns > 0 && (rules.MercyAppliesLastInning || !rules.MercyAppliesLastInning && game.State.Inning == rules.MaxInnings) {
		if game.State.HomeScore-game.State.AwayScore >= rules.GameMercyRuns || game.State.AwayScore-game.State.HomeScore >= rules.GameMercyRuns {
			// For now, we simply leave scores as-is; higher-level callers may treat games past mercy as finished.
		}
	}

	return game, nil
}
