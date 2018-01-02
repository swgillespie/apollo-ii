package engine

import (
	"fmt"
)

// A Move is a transformation on the game board. Each player can make moves
// to advance the game. Moves are all encoded in a compact manner as 16-bit
// integers for the sake of performance.
//
// ## Encoding
// The encoding of a move is like this:
//  * 6 bits - source square
//  * 6 bits - destination square
//  * 1 bit  - promotion bit
//  * 1 bit  - capture bit
//  * 1 bit  - "special 0" square
//  * 1 bit  - "special 1" square
// The "special" bits are overloaded, because chess has a
// number of "special" moves that do not fit nicely into
// a compact representation. Here is a full table of
// the encoding strategy:
//
// | Promo | Capt  | Spc 0 | Spc 1 | Move                   |
// |-------|-------|-------|-------|------------------------|
// | 0     | 0     | 0     | 0     | Quiet                  |
// | 0     | 0     | 0     | 1     | Double Pawn            |
// | 0     | 0     | 1     | 0     | King Castle            |
// | 0     | 0     | 1     | 1     | Queen Castle           |
// | 0     | 1     | 0     | 0     | Capture                |
// | 0     | 1     | 0     | 1     | En Passant Capture     |
// | 1     | 0     | 0     | 0     | Knight Promote         |
// | 1     | 0     | 0     | 1     | Bishop Promote         |
// | 1     | 0     | 1     | 0     | Rook Promote           |
// | 1     | 0     | 1     | 1     | Queen Promote          |
// | 1     | 1     | 0     | 0     | Knight Promote Capture |
// | 1     | 1     | 0     | 1     | Bishop Promote Capture |
// | 1     | 1     | 1     | 0     | Rook Promote Capture   |
// | 1     | 1     | 1     | 1     | Queen Promote Capture  |
//
// Thanks to https://chessprogramming.wikispaces.com/Encoding+Moves
// for the details.
type Move uint16

const (
	sourceMask      = 0xFC00
	destinationMask = 0x03F0
	promoBit        = 0x0008
	captureBit      = 0x0004
	special0Bit     = 0x0002
	special1Bit     = 0x0001
	attrMask        = 0x000F
)

// MakeQuietMove constructs a new quiet move from the source
// square to the destination square.
func MakeQuietMove(source, dest Square) Move {
	sourceBits := uint16(source) << 10
	destBits := uint16(dest) << 4
	return Move(sourceBits | destBits)
}

// MakeCaptureMove constructs a new capture move from the source square
// to the destination square.
func MakeCaptureMove(source, dest Square) Move {
	mov := MakeQuietMove(source, dest)
	mov |= captureBit
	return mov
}

func MakeEnPassantMove(source, dest Square) Move {
	mov := MakeCaptureMove(source, dest)
	mov |= special1Bit
	return mov
}

func MakeDoublePawnPushMove(source, dest Square) Move {
	mov := MakeQuietMove(source, dest)
	mov |= special1Bit
	return mov
}

func MakePromotionMove(source, dest Square, piece PieceKind) Move {
	mov := MakeQuietMove(source, dest)
	mov |= promoBit
	switch piece {
	case Knight:
		mov |= 0
	case Bishop:
		mov |= 1
	case Rook:
		mov |= 2
	case Queen:
		mov |= 3
	}

	return mov
}

func MakePromotionCaptureMove(source, dest Square, piece PieceKind) Move {
	mov := MakePromotionMove(source, dest, piece)
	mov |= captureBit
	return mov
}

func MakeKingsideCastleMove(source, dest Square) Move {
	mov := MakeQuietMove(source, dest)
	mov |= special0Bit
	return mov
}

func MakeQueensideCastleMove(source, dest Square) Move {
	mov := MakeQuietMove(source, dest)
	mov |= special0Bit | special1Bit
	return mov
}

func MakeNullMove(source, dest Square) Move {
	return Move(0)
}

func (m Move) Source() Square {
	return Square((m & sourceMask) >> 10)
}

func (m Move) Destination() Square {
	return Square((m & destinationMask) >> 4)
}

func (m Move) IsQuiet() bool {
	return (m & attrMask) == 0
}

func (m Move) IsCapture() bool {
	return (m & captureBit) != 0
}

func (m Move) IsEnPassant() bool {
	return (m & attrMask) == 5
}

func (m Move) IsDoublePawnPush() bool {
	return (m & attrMask) == 1
}

func (m Move) IsPromotion() bool {
	return (m & promoBit) == promoBit
}

func (m Move) IsKingsideCastle() bool {
	return (m & attrMask) == 2
}

func (m Move) IsQueensideCastle() bool {
	return (m & attrMask) == 3
}

func (m Move) IsCastle() bool {
	return m.IsKingsideCastle() || m.IsQueensideCastle()
}

func (m Move) IsNull() bool {
	return m == 0
}

func (m Move) PromotionPiece() PieceKind {
	if !m.IsPromotion() {
		panic("PromotionPiece called on non-promotion move")
	}

	piece := m & (special0Bit | special1Bit)
	switch piece {
	case 0:
		return Knight
	case 1:
		return Bishop
	case 2:
		return Rook
	case 3:
		return Queen
	}

	// piece is 2 bits wide by definition. It is not possible for it to have
	// any values other than 0, 1, 2, and 3.
	panic("unreachable code in PromotionPiece")
}

// UciString returns a UCI-encoded representation of this move.
func (m Move) UciString() string {
	if !m.IsPromotion() {
		return fmt.Sprintf("%s%s", m.Source(), m.Destination())
	}

	return fmt.Sprintf("%s%s%s", m.Source(), m.Destination(), m.PromotionPiece())
}

func (m Move) String() string {
	return m.UciString()
}
