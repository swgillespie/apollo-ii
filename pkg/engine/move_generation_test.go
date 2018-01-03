package engine

import "testing"
import "github.com/stretchr/testify/assert"

// Targeted tests for the move generator.
//
// PERFT is an insanely good way to test the move generator, but
// these tests are much more targeted in that they test specific
// situations. Failed PERFT tests should usually result in new
// movegen test cases.
//
// Pretty much every test from this file came out of an investigation
// as to why PERFT numbers weren't correct.
func AssertHasMove(t *testing.T, fen string, mov Move) {
	pos, err := MakePositionFromFen(fen)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	moves := pos.PseudolegalMoves()
	for _, generated := range moves {
		if generated == mov {
			return
		}
	}

	assert.FailNow(t, "move-check", "expected to generate move `%s` but did not, generated %v", mov.String(), moves)
}

func TestMoveGeneration(t *testing.T) {
	Initialize()
	t.Parallel()
	t.Run("early-game-rook", func(tt *testing.T) {
		// both black and white have bumped their a-rank pawns. Now white
		// can move their rook to a1 to a2.
		AssertHasMove(tt, "rnbqkbnr/1ppppppp/p7/8/8/P7/1PPPPPPP/RNBQKBNR w KQkq - 0 1", MakeQuietMove(A1, A2))
	})

	t.Run("early-game-king", func(tt *testing.T) {
		AssertHasMove(tt, "rnbqkbnr/1ppppppp/p7/8/8/3P4/PPP1PPPP/RNBQKBNR w KQkq - 0 2", MakeQuietMove(E1, D2))
	})
}
