package services

import (
	"strings"

	"github.com/ajarvis3/kickball-go/internal/data/domain"
	"github.com/ajarvis3/kickball-go/internal/types"
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

	// Game initialization (inning/half/runs) is handled at game creation.
	// ApplyAtBat assumes game.State has been initialized by the creator.

	res := strings.ToLower(strings.TrimSpace(atbat.Result))

	// map inning/half -> index: (inning-1)*2 + (0 for top, 1 for bottom)
	idx := (game.State.Inning - 1) * 2
	if game.State.Half == types.HalfBottom {
		idx += 1
	}

	// Apply runs (RBI) for the at-bat
	applyRBI(&game, atbat, idx)

	// Apply outs for the at-bat
	applyOuts(&game, res)

	// Advance inning/half when outs reach 3 or more
	for game.State.Outs >= 3 {
		advanceInning(&game)
	}

	// Innings mercy rule check (rules apply AND inning runs exceed mercy threshold)
	if e.Rules.DoesInningMercyApply(rules, game, idx) {
		adjustedRuns := game.State.InningRuns[idx] - rules.MercyRunsPerInning
		if atbat.TeamID == game.HomeTeamID {
			game.State.HomeScore -= adjustedRuns
		} else {
			game.State.AwayScore -= adjustedRuns
		}
		game.State.InningRuns[idx] = rules.MercyRunsPerInning
		// reset outs and advance half/inning
		advanceInning(&game)
	}

	// Basic mercy rule check (game-level)
	if rules.GameMercyRuns > 0 {
		if game.State.HomeScore-game.State.AwayScore >= rules.GameMercyRuns || game.State.AwayScore-game.State.HomeScore >= rules.GameMercyRuns {
			game.State.Inning = rules.MaxInnings + 1
			game.State.Outs = 0
			game.State.Half = types.Half("")
		}
	}

	return game, nil
}

// advanceInning resets outs and advances the half/inning on the provided game.
// It's unexported and intended for use only within this file.
func advanceInning(g *domain.Game) {
	if g == nil {
		return
	}
	g.State.Outs = 0
	if g.State.Half == types.HalfTop {
		g.State.Half = types.HalfBottom
	} else {
		g.State.Half = types.HalfTop
		g.State.Inning += 1
	}
}

// applyRBI updates scores and inning runs for at-bats that produce RBI.
// It is unexported and used only within this file.
func applyRBI(g *domain.Game, atbat domain.AtBat, idx int) {
	res := strings.ToLower(strings.TrimSpace(atbat.Result))
	if isRBIResult(res) && atbat.RBI > 0 {
		if atbat.TeamID == g.HomeTeamID {
			g.State.HomeScore += atbat.RBI
		} else {
			g.State.AwayScore += atbat.RBI
		}
		if idx >= 0 && idx < len(g.State.InningRuns) {
			g.State.InningRuns[idx] += atbat.RBI
		}
	}
}

var rbiResults = map[string]struct{}{
	"single":         {},
	"double":         {},
	"triple":         {},
	"homerun":        {},
	"walk":           {},
	"sacrifice":      {},
	"error":          {},
	"fielderschoice": {},
}

func isRBIResult(res string) bool {
	_, ok := rbiResults[res]
	return ok
}

// doesInningMercyApply returns true when the mercy rule should be checked
// for the current inning according to the league rules.
// (was moved to RulesEngine)

// applyOuts increments the game's outs according to result
func applyOuts(g *domain.Game, res string) {
	switch res {
	case "out", "strikeout":
		g.State.Outs += 1
	case "doubleplay":
		g.State.Outs += 2
	case "tripleplay":
		g.State.Outs += 3
	}
}
