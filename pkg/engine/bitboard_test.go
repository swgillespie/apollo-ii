package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitboards(t *testing.T) {
	t.Parallel()
	t.Run("smoke", func(tt *testing.T) {
		empty := EmptyBitboard
		for sq := A1; sq <= H8; sq++ {
			assert.False(tt, empty.Test(sq))
		}
	})

	t.Run("set", func(tt *testing.T) {
		board := EmptyBitboard
		board.Set(A4)
		assert.True(tt, board.Test(A4))
		for sq := A1; sq <= H8; sq++ {
			if sq != A4 {
				assert.False(tt, board.Test(sq))
			}
		}
	})

	t.Run("unset", func(tt *testing.T) {
		board := FullBitboard
		board.Unset(H3)
		assert.False(tt, board.Test(H3))
		for sq := A1; sq <= H8; sq++ {
			if sq != H3 {
				assert.True(tt, board.Test(sq))
			}
		}
	})

	t.Run("union", func(tt *testing.T) {
		board := EmptyBitboard
		board.Set(A1)
		other := EmptyBitboard
		other.Set(B5)
		assert.True(tt, board.Test(A1))
		assert.False(tt, board.Test(B5))
		assert.True(tt, other.Test(B5))
		assert.False(tt, other.Test(A1))
		union := board | other
		assert.True(tt, union.Test(A1))
		assert.True(tt, union.Test(B5))
	})

	t.Run("intersection", func(tt *testing.T) {
		board := EmptyBitboard
		other := EmptyBitboard
		board.Set(A1)
		board.Set(A2)
		other.Set(B1)
		other.Set(A2)
		intersect := board & other
		assert.False(tt, intersect.Test(A1))
		assert.True(tt, intersect.Test(A2))
		assert.False(tt, intersect.Test(B1))
	})

	t.Run("empty-iter", func(tt *testing.T) {
		board := EmptyBitboard
		var squares []Square
		iter := board.Iter()
		for sq, ok := iter.Next(); ok; sq, ok = iter.Next() {
			squares = append(squares, sq)
		}

		assert.Len(tt, squares, 0)
	})

	t.Run("iter-one", func(tt *testing.T) {
		board := EmptyBitboard
		board.Set(A4)
		var squares []Square
		iter := board.Iter()
		for sq, ok := iter.Next(); ok; sq, ok = iter.Next() {
			squares = append(squares, sq)
		}

		assert.Len(tt, squares, 1)
		assert.Equal(tt, squares[0], A4)
	})

	t.Run("iter-many", func(tt *testing.T) {
		board := FullBitboard
		var squares []Square

		iter := board.Iter()
		for sq, ok := iter.Next(); ok; sq, ok = iter.Next() {
			squares = append(squares, sq)
		}

		assert.Len(tt, squares, 64)
		for sq := A1; sq <= H8; sq++ {
			assert.Contains(tt, squares, sq)
		}
	})

	t.Run("rank", func(tt *testing.T) {
		board := EmptyBitboard
		board.Set(E4)
		board.Set(F4)
		board.Set(F5)
		board = board.Rank(Rank4)
		assert.True(tt, board.Test(E4))
		assert.True(tt, board.Test(F4))
		assert.False(tt, board.Test(F5))
	})

	t.Run("file", func(tt *testing.T) {
		board := EmptyBitboard
		board.Set(E4)
		board.Set(F4)
		board.Set(F5)
		board = board.File(FileF)
		assert.False(tt, board.Test(E4))
		assert.True(tt, board.Test(F4))
		assert.True(tt, board.Test(F5))
	})
}

// Bitboard square iteration is performance critical. It needs to be fast
// and zero-alloc.
//
// It would be super cool if I could get Go to fail the test run if the
// alloc count is not zero.
func BenchmarkBitboardIter(b *testing.B) {
	b.Run("full-board", func(bb *testing.B) {
		bb.ReportAllocs()
		board := FullBitboard
		for i := 0; i < bb.N; i++ {
			iter := board.Iter()
			for _, ok := iter.Next(); ok; _, ok = iter.Next() {
			}
		}
	})

	b.Run("10-element-board", func(bb *testing.B) {
		bb.ReportAllocs()
		board := EmptyBitboard
		board.Set(A1)
		board.Set(H4)
		board.Set(B3)
		board.Set(B1)
		board.Set(C3)
		board.Set(D1)
		board.Set(D4)
		board.Set(D7)
		board.Set(H8)
		board.Set(A8)
		b.ResetTimer()
		for i := 0; i < bb.N; i++ {
			iter := board.Iter()
			for _, ok := iter.Next(); ok; _, ok = iter.Next() {
			}
		}
	})

	b.Run("empty-board", func(bb *testing.B) {
		bb.ReportAllocs()
		board := EmptyBitboard
		for i := 0; i < bb.N; i++ {
			iter := board.Iter()
			for _, ok := iter.Next(); ok; _, ok = iter.Next() {
			}
		}
	})
}
