package perft

import (
	"fmt"

	"github.com/swgillespie/apollo-ii/pkg/engine"
)

type Results struct {
	Nodes      uint64
	Captures   uint64
	EnPassants uint64
	Castles    uint64
	Promotions uint64
	Checks     uint64
	Checkmates uint64
}

func Perft(fen string, depth int) (*Results, error) {
	if depth < 0 {
		return nil, fmt.Errorf("invalid ply depth: %d", depth)
	}

	pos, err := engine.MakePositionFromFen(fen)
	if err != nil {
		return nil, err
	}

	results := new(Results)
	perftImpl(results, pos, depth)
	return results, nil
}

func perftImpl(results *Results, pos *engine.Position, depth int) {
	if depth == 0 {
		results.Nodes++
		return
	}

	seenLegalMove := false
	for _, move := range pos.PseudolegalMoves() {
		newPos := pos.Clone()
		toMove := pos.SideToMove()
		newPos.ApplyMove(move)
		if !newPos.IsCheck(toMove) {
			seenLegalMove = true
			if move.IsCapture() {
				results.Captures++
			}

			if move.IsEnPassant() {
				results.EnPassants++
			}

			if move.IsKingsideCastle() || move.IsQueensideCastle() {
				results.Castles++
			}

			if move.IsPromotion() {
				results.Promotions++
			}

			if newPos.IsCheck(toMove.Toggle()) {
				results.Checks++
			}

			perftImpl(results, newPos, depth-1)
		}
	}

	if !seenLegalMove {
		results.Checkmates++
	}
}
