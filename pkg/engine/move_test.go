package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoveEncodings(t *testing.T) {
	t.Parallel()
	t.Run("quiet", func(tt *testing.T) {
		mov := MakeQuietMove(A4, A5)
		assert.True(tt, mov.IsQuiet())
		assert.Equal(tt, mov.Source(), A4)
		assert.Equal(tt, mov.Destination(), A5)
	})

	t.Run("capture", func(tt *testing.T) {
		mov := MakeCaptureMove(B4, C4)
		assert.True(tt, mov.IsCapture())
		assert.Equal(tt, mov.Source(), B4)
		assert.Equal(tt, mov.Destination(), C4)
	})

	t.Run("en-passant", func(tt *testing.T) {
		mov := MakeEnPassantMove(C4, E4)
		assert.True(tt, mov.IsCapture())
		assert.True(tt, mov.IsEnPassant())
		assert.False(tt, mov.IsQuiet())
	})

	t.Run("double-pawn-push", func(tt *testing.T) {
		mov := MakeDoublePawnPushMove(E2, E4)
		assert.True(tt, mov.IsDoublePawnPush())
		assert.False(tt, mov.IsCapture())
		assert.False(tt, mov.IsQuiet())
	})

	for kind := Knight; kind != King; kind++ {
		t.Run(fmt.Sprintf("promotion-%s", kind), func(tt *testing.T) {
			mov := MakePromotionMove(H7, H8, kind)
			assert.True(tt, mov.IsPromotion())
			assert.False(tt, mov.IsCapture())
			assert.Equal(tt, kind, mov.PromotionPiece())
		})
	}
}

func TestMoveUciStrings(t *testing.T) {
	t.Parallel()
	t.Run("quiet", func(tt *testing.T) {
		mov := MakeQuietMove(E4, E5)
		assert.Equal(tt, "e4e5", mov.UciString())
	})

	t.Run("promotion", func(tt *testing.T) {
		mov := MakePromotionMove(H7, H8, Queen)
		assert.Equal(tt, "h7h8q", mov.UciString())
	})
}
