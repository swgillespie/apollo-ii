package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFen(t *testing.T) {
	t.Parallel()
	t.Run("smoke-empty-board", func(tt *testing.T) {
		pos, err := MakePositionFromFen("8/8/8/8/8/8/8/8 w - - 0 0")
		// these blow up the test if they fail
		if !assert.Nil(tt, err) {
			tt.FailNow()
		}

		if !assert.NotNil(tt, pos) {
			tt.FailNow()
		}

		// empty board, white's turn to move.
		all := pos.White() | pos.Black()
		assert.Zero(tt, all.Count())
		assert.Equal(tt, White, pos.SideToMove())

		// no EP square
		assert.False(tt, pos.HasEnPassantSquare())

		// nobody can castle
		assert.False(tt, pos.CanCastleKingside(White))
		assert.False(tt, pos.CanCastleKingside(Black))
		assert.False(tt, pos.CanCastleQueenside(White))
		assert.False(tt, pos.CanCastleQueenside(White))

		// clocks are zero
		assert.Zero(tt, pos.HalfmoveClock())
		assert.Zero(tt, pos.FullmoveClock())
	})

	t.Run("smoke-test-starting-position", func(tt *testing.T) {
		// this is the FEN representation of the starting configuration
		// of a standard game of chess
		pos, err := MakePositionFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		if !assert.Nil(tt, err) {
			tt.FailNow()
		}

		if !assert.NotNil(tt, pos) {
			tt.FailNow()
		}

		checkSquare := func(squareStr string, kind PieceKind, color Color) {
			square, err := MakeSquareFromString(squareStr)
			if !assert.Nil(tt, err) {
				tt.Fatalf("unexpected error: %s", err.Error())
			}

			piece, ok := pos.PieceAt(square)
			if !assert.True(tt, ok) {
				return
			}

			assert.Equal(tt, kind, piece.kind, "square `%s` had wrong piece kind", squareStr)
			assert.Equal(tt, color, piece.color, "square `%s` had wrong piece color", squareStr)
		}

		checkVacant := func(square Square) {
			_, ok := pos.PieceAt(square)
			assert.False(tt, ok, "square `%s` was not vacant", square.String())
		}

		checkSquare("a1", Rook, White)
		checkSquare("b1", Knight, White)
		checkSquare("c1", Bishop, White)
		checkSquare("d1", Queen, White)
		checkSquare("e1", King, White)
		checkSquare("f1", Bishop, White)
		checkSquare("g1", Knight, White)
		checkSquare("h1", Rook, White)

		checkSquare("a2", Pawn, White)
		checkSquare("b2", Pawn, White)
		checkSquare("c2", Pawn, White)
		checkSquare("d2", Pawn, White)
		checkSquare("e2", Pawn, White)
		checkSquare("f2", Pawn, White)
		checkSquare("g2", Pawn, White)
		checkSquare("h2", Pawn, White)

		for sq := A3; sq < A7; sq++ {
			checkVacant(sq)
		}

		checkSquare("a7", Pawn, Black)
		checkSquare("b7", Pawn, Black)
		checkSquare("c7", Pawn, Black)
		checkSquare("d7", Pawn, Black)
		checkSquare("e7", Pawn, Black)
		checkSquare("f7", Pawn, Black)
		checkSquare("g7", Pawn, Black)
		checkSquare("h7", Pawn, Black)

		checkSquare("a8", Rook, Black)
		checkSquare("b8", Knight, Black)
		checkSquare("c8", Bishop, Black)
		checkSquare("d8", Queen, Black)
		checkSquare("e8", King, Black)
		checkSquare("f8", Bishop, Black)
		checkSquare("g8", Knight, Black)
		checkSquare("h8", Rook, Black)

		// no EP square
		assert.False(tt, pos.HasEnPassantSquare())

		// everyone can castle in every way
		assert.True(tt, pos.CanCastleKingside(White))
		assert.True(tt, pos.CanCastleKingside(Black))
		assert.True(tt, pos.CanCastleQueenside(White))
		assert.True(tt, pos.CanCastleQueenside(Black))

		// halfmove clock is zero
		assert.Zero(tt, pos.HalfmoveClock())

		// fullmove clock is one
		assert.Equal(tt, uint32(1), pos.FullmoveClock())
	})

	t.Run("smoke-en-passant", func(tt *testing.T) {
		pos, _ := MakePositionFromFen("8/8/8/8/8/8/8/8 w - e3 0 0")
		if !assert.NotNil(tt, pos) {
			tt.FailNow()
		}

		assert.True(tt, pos.HasEnPassantSquare())
		assert.Equal(tt, E3, pos.EnPassantSquare())
	})

	t.Run("black-side-to-move", func(tt *testing.T) {
		pos, _ := MakePositionFromFen("8/8/8/8/8/8/8/8 b - - 0 0")
		if !assert.NotNil(tt, pos) {
			tt.FailNow()
		}

		assert.Equal(tt, Black, pos.SideToMove())
	})

	t.Run("empty", func(tt *testing.T) {
		_, err := MakePositionFromFen("")
		if assert.Error(tt, err) {
			assert.Equal(tt, FenUnknownRuneOrEofError, err)
		}
	})

	t.Run("unknown-piece", func(tt *testing.T) {
		_, err := MakePositionFromFen("z7/8/8/8/8/8/8/8 w - - 0 0")
		if assert.Error(tt, err) {
			assert.Equal(tt, FenUnknownRuneOrEofError, err)
		}
	})

	t.Run("invalid-digit", func(tt *testing.T) {
		_, err := MakePositionFromFen("9/8/8/8/8/8/8/8 w - - 0 0")
		if assert.Error(tt, err) {
			assert.Equal(tt, FenInvalidDigitError, err)
		}
	})

	t.Run("not-sum-to-8", func(tt *testing.T) {
		_, err := MakePositionFromFen("pppp5/8/8/8/8/8/8/8 w - - 0 0")
		if assert.Error(tt, err) {
			assert.Equal(tt, FenSumToEightError, err)
		}
	})

	t.Run("bad-side-to-move", func(tt *testing.T) {
		_, err := MakePositionFromFen("8/8/8/8/8/8/8/8 c - - 0 0")
		if assert.Error(tt, err) {
			assert.Equal(tt, FenInvalidSideToMoveError, err)
		}
	})

	t.Run("bad-castle-status", func(tt *testing.T) {
		_, err := MakePositionFromFen("8/8/8/8/8/8/8/8 w a - 0 0")
		if assert.Error(tt, err) {
			assert.Equal(tt, FenInvalidCastleStatusError, err)
		}
	})

	t.Run("bad-en-passant", func(tt *testing.T) {
		_, err := MakePositionFromFen("8/8/8/8/8/8/8/8 w - 88 0 0")
		if assert.Error(tt, err) {
			assert.Equal(tt, FenInvalidEnPassantError, err)
		}
	})

	t.Run("empty-halfmove", func(tt *testing.T) {
		_, err := MakePositionFromFen("8/8/8/8/8/8/8/8 w - - q 0")
		if assert.Error(tt, err) {
			assert.Equal(tt, FenInvalidHalfmoveError, err)
		}
	})

	t.Run("invalid-halfmove", func(tt *testing.T) {
		_, err := MakePositionFromFen("8/8/8/8/8/8/8/8 w - - 4294967296 0")
		if assert.Error(tt, err) {
			assert.Equal(tt, FenInvalidHalfmoveError, err)
		}
	})

	t.Run("empty-fullmove", func(tt *testing.T) {
		_, err := MakePositionFromFen("8/8/8/8/8/8/8/8 w - - 0 q")
		if assert.Error(tt, err) {
			assert.Equal(tt, FenInvalidFullmoveError, err)
		}
	})

	t.Run("fullmove-early-end", func(tt *testing.T) {
		_, err := MakePositionFromFen("8/8/8/8/8/8/8/8 w - - 0")
		if assert.Error(tt, err) {
			assert.Equal(tt, FenEndOfFileError, err)
		}
	})

	t.Run("invalid-fullmove", func(tt *testing.T) {
		_, err := MakePositionFromFen("8/8/8/8/8/8/8/8 w - - 0 4294967296")
		if assert.Error(tt, err) {
			assert.Equal(tt, FenInvalidFullmoveError, err)
		}
	})

	t.Run("king-queen-castle-status", func(tt *testing.T) {
		pos, err := MakePositionFromFen("8/8/8/8/8/8/8/4K2R w KQ - 0 1")
		if !assert.NoError(tt, err) {
			tt.FailNow()
		}

		assert.True(tt, pos.CanCastleKingside(White))
		assert.True(tt, pos.CanCastleQueenside(White))
		assert.False(tt, pos.CanCastleKingside(Black))
		assert.False(tt, pos.CanCastleQueenside(Black))
	})
}

var fenRoundtripTests = [...]string{
	"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	"8/6Q1/5b2/4r3/3k4/2R5/1N6/P7 w - - 0 1",
}

func TestFenRoundTrip(t *testing.T) {
	t.Parallel()
	for _, fen := range fenRoundtripTests {
		t.Run(fmt.Sprintf("fen-roundtrip-%s", fen), func(tt *testing.T) {
			pos, err := MakePositionFromFen(fen)
			if !assert.NoError(tt, err) {
				tt.FailNow()
			}

			assert.Equal(tt, fen, pos.AsFen())
		})
	}
}
