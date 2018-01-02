package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoveApplication(t *testing.T) {
	t.Parallel()
	t.Run("smoke-test-opening-pawn", func(tt *testing.T) {
		pos := MakeDefaultPosition()

		// nothing fancy, move a pawn up one.
		pos.ApplyMove(MakeQuietMove(E2, E3))

		// it should now be Black's turn to move.
		assert.Equal(tt, Black, pos.SideToMove())

		// the fullmove clock shouldn't have incremented
		// (it only increments every Black move)
		assert.Equal(tt, uint32(1), pos.FullmoveClock())

		// a pawn moved, so the halfmove clock should be zero.
		assert.Equal(tt, uint32(0), pos.HalfmoveClock())

		// there should be a pawn on e3
		piece, ok := pos.PieceAt(E3)
		if !assert.True(tt, ok) {
			tt.FailNow()
		}

		assert.Equal(tt, Pawn, piece.kind)
		assert.Equal(tt, White, piece.color)

		// there should not be a pawn on e2
		_, ok = pos.PieceAt(E2)
		assert.False(tt, ok)
	})

	t.Run("en-passant-reset", func(tt *testing.T) {
		// EP square at e3, black to move
		pos, err := MakePositionFromFen("8/8/8/8/4Pp2/8/8/8 b - e3 0 1")
		if !assert.Nil(tt, err) {
			tt.FailNow()
		}

		// black not taking EP opportunity
		pos.ApplyMove(MakeQuietMove(F4, F3))

		// EP no longer possible.
		assert.Equal(tt, White, pos.SideToMove())
		assert.False(tt, pos.HasEnPassantSquare())
	})

	t.Run("double-pawn-push-sets-ep", func(tt *testing.T) {
		// white to move
		pos, err := MakePositionFromFen("8/8/8/8/8/8/4P3/8 w - - 0 1")
		if !assert.Nil(tt, err) {
			tt.FailNow()
		}

		// white double-pawn pushes
		pos.ApplyMove(MakeDoublePawnPushMove(E2, E4))

		// now black to move, with EP square set
		assert.Equal(tt, Black, pos.SideToMove())
		assert.True(tt, pos.HasEnPassantSquare())
		assert.Equal(tt, E3, pos.EnPassantSquare())
	})

	t.Run("basic-capture", func(tt *testing.T) {
		pos, err := MakePositionFromFen("8/8/8/8/5p2/4P3/8/8 w - - 2 1")
		if !assert.Nil(tt, err) {
			tt.FailNow()
		}

		pos.ApplyMove(MakeCaptureMove(E3, F4))

		// There should be a white pawn on F4
		piece, ok := pos.PieceAt(F4)
		if !assert.True(tt, ok) {
			tt.FailNow()
		}

		assert.Equal(tt, Pawn, piece.kind)
		assert.Equal(tt, White, piece.color)

		// There should be no piece on E3
		_, ok = pos.PieceAt(E3)
		assert.False(tt, ok)

		// The halfmove clock should reset (capture)
		assert.Equal(tt, uint32(0), pos.HalfmoveClock())
	})

	t.Run("non-pawn-quiet-move", func(tt *testing.T) {
		pos, err := MakePositionFromFen("8/8/8/8/8/8/4B3/8 w - - 5 2")
		if !assert.Nil(tt, err) {
			tt.FailNow()
		}

		pos.ApplyMove(MakeQuietMove(E2, G4))

		// the halfmove clock should not be reset.
		assert.Equal(tt, uint32(6), pos.HalfmoveClock())
	})

	t.Run("moving-king-castle-status", func(tt *testing.T) {
		pos, err := MakePositionFromFen("8/8/8/8/8/8/8/4K2R w KQ - 0 1")
		if !assert.Nil(tt, err) {
			tt.FailNow()
		}

		// white's turn to move, white moves its king.
		pos.ApplyMove(MakeQuietMove(E1, E2))

		// white can't castle anymore.
		assert.False(tt, pos.CanCastleKingside(White))
		assert.False(tt, pos.CanCastleQueenside(White))
	})

	t.Run("moving-kingside-rook-castle-status", func(tt *testing.T) {
		pos, err := MakePositionFromFen("8/8/8/8/8/8/8/4K2R w KQ - 0 1")
		if !assert.NoError(tt, err) {
			tt.FailNow()
		}

		// white's turn to move, white moves its kingside rook.
		pos.ApplyMove(MakeQuietMove(H1, G1))

		// white can't castle kingside anymore
		assert.False(tt, pos.CanCastleKingside(White))
		assert.True(tt, pos.CanCastleQueenside(White))
	})

	t.Run("moving-queenside-rook-castle-status", func(tt *testing.T) {
		pos, err := MakePositionFromFen("8/8/8/8/8/8/8/R3K3 w KQ - 0 1")
		if !assert.NoError(tt, err) {
			tt.FailNow()
		}

		// white's turn to move, white moves its queenside rook.
		pos.ApplyMove(MakeQuietMove(A1, B1))

		// white can't castle queenside anymore
		assert.True(tt, pos.CanCastleKingside(White))
		assert.False(tt, pos.CanCastleQueenside(White))
	})

	t.Run("rook-capture-castle-status", func(tt *testing.T) {
		pos, err := MakePositionFromFen("8/8/8/8/8/7r/4P3/R3K2R b KQ - 0 1")
		if !assert.NoError(tt, err) {
			tt.FailNow()
		}

		// tests that we can't capture if there's no rook on the target
		// square, even if the rooks themselves never moved (i.e. they
		// were captured on their starting square)

		// black to move, black captures the rook at H1
		pos.ApplyMove(MakeCaptureMove(H3, H1))

		// white to move, white pushes the pawn
		pos.ApplyMove(MakeDoublePawnPushMove(E2, E4))

		// black to move, black moves the rook
		pos.ApplyMove(MakeQuietMove(H1, H5))

		// white moves the queenside rook to the kingside rook
		// start location
		pos.ApplyMove(MakeQuietMove(A1, A2))
		pos.ApplyMove(MakeQuietMove(H5, H6))
		pos.ApplyMove(MakeQuietMove(A2, H2))
		pos.ApplyMove(MakeQuietMove(H6, H7))
		pos.ApplyMove(MakeQuietMove(H2, H1))

		// white shouldn't be able to castle kingside, despite
		// there being a rook on the kingside rook square
		// and us never moving the kingside rook
		assert.False(tt, pos.CanCastleKingside(White))
	})

	t.Run("en-passant-capture", func(tt *testing.T) {
		// tests that we remove an ep-captured piece from its
		// actual location and not try to remove the EP-square
		pos, err := MakePositionFromFen("8/8/8/3pP3/8/8/8/8 w - d6 0 1")
		if !assert.NoError(tt, err) {
			tt.FailNow()
		}

		// white to move, white EP-captures the pawn
		pos.ApplyMove(MakeEnPassantMove(E5, D6))

		// there should not be a piece at D5 anymore
		_, ok := pos.PieceAt(D5)
		assert.False(tt, ok)

		// the white pawn should be at the EP-square
		whitePawn, ok := pos.PieceAt(D6)
		if !assert.True(tt, ok) {
			tt.FailNow()
		}

		assert.Equal(tt, White, whitePawn.color)
		assert.Equal(tt, Pawn, whitePawn.kind)
	})

	t.Run("basic-promotion", func(tt *testing.T) {
		pos, err := MakePositionFromFen("8/4P3/8/8/8/8/8/8 w - - 0 1")
		if !assert.NoError(tt, err) {
			tt.FailNow()
		}

		// white to move, white promotes the pawn on e7
		pos.ApplyMove(MakePromotionMove(E7, E8, Queen))

		// there should be a queen on e8
		queen, ok := pos.PieceAt(E8)
		if !assert.True(tt, ok) {
			tt.FailNow()
		}

		assert.Equal(tt, White, queen.color)
		assert.Equal(tt, Queen, queen.kind)
	})

	t.Run("queenside-castle", func(tt *testing.T) {
		pos, err := MakePositionFromFen("8/8/8/8/8/8/8/R3K3 w Q - 0 1")
		if !assert.NoError(tt, err) {
			tt.FailNow()
		}

		// white to move, white castles queenside
		pos.ApplyMove(MakeQueensideCastleMove(E1, C1))

		rook, ok := pos.PieceAt(D1)
		if !assert.True(tt, ok) {
			tt.FailNow()
		}

		assert.Equal(tt, White, rook.color)
		assert.Equal(tt, Rook, rook.kind)

		king, ok := pos.PieceAt(C1)
		if !assert.True(tt, ok) {
			tt.FailNow()
		}

		assert.Equal(tt, White, king.color)
		assert.Equal(tt, King, king.kind)
	})

	t.Run("kingside-castle", func(tt *testing.T) {
		pos, err := MakePositionFromFen("8/8/8/8/8/8/8/4K2R w K - 0 1")
		if !assert.NoError(tt, err) {
			tt.FailNow()
		}

		// white to move, white castles kingside
		pos.ApplyMove(MakeKingsideCastleMove(E1, G1))

		rook, ok := pos.PieceAt(F1)
		if !assert.True(tt, ok) {
			tt.FailNow()
		}

		assert.Equal(tt, White, rook.color)
		assert.Equal(tt, Rook, rook.kind)

		king, ok := pos.PieceAt(G1)
		if !assert.True(tt, ok) {
			tt.FailNow()
		}

		assert.Equal(tt, White, king.color)
		assert.Equal(tt, King, king.kind)
	})
}
