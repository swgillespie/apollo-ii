package engine

import (
	"math/bits"
)

// A Bitboard is a 64-bit integer which one bit represents one of the
// eight squares on the board. Bitboards are used in a variety of scenarios
// to represent the board itself and the pieces upon it.
type Bitboard uint64

const (
	// EmptyBitboard is the empty bitboard, analagous to the empty set.
	EmptyBitboard = Bitboard(0)

	// FullBitboard is the full bitboard, the complement of the empty set.
	FullBitboard = Bitboard(0xFFFFFFFFFFFFFFFF)
)

// Masks for masking off particular ranks or files, indexed by rank or file.
var rankMasks = [...]uint64{
	0x00000000000000FF,
	0x000000000000FF00,
	0x0000000000FF0000,
	0x00000000FF000000,
	0x000000FF00000000,
	0x0000FF0000000000,
	0x00FF000000000000,
	0xFF00000000000000}

var fileMasks = [...]uint64{
	0x0101010101010101,
	0x0202020202020202,
	0x0404040404040404,
	0x0808080808080808,
	0x1010101010101010,
	0x2020202020202020,
	0x4040404040404040,
	0x8080808080808080}

// Test tests whether or not a square is a member of this bitboard.
func (b Bitboard) Test(square Square) bool {
	return (uint64(b) & (uint64(1) << square)) != 0
}

func (b Bitboard) Rank(rank Rank) Bitboard {
	return Bitboard(uint64(b) & rankMasks[rank])
}

func (b Bitboard) File(file File) Bitboard {
	return Bitboard(uint64(b) & fileMasks[file])
}

// Set sets a square to be a member of this bitboard.
func (b *Bitboard) Set(square Square) {
	*b = Bitboard(uint64(*b) | (uint64(1) << square))
}

func (b *Bitboard) Unset(square Square) {
	*b = Bitboard(uint64(*b) & ^(uint64(1) << square))
}

// Count returns the number of squares in this bitboard.
func (b Bitboard) Count() int {
	return bits.OnesCount64(uint64(b))
}

// Empty returns whether or not this bitboard is empty.
func (b Bitboard) Empty() bool {
	return b == 0
}

func (b Bitboard) Iter() BitboardIterator {
	return BitboardIterator{b}
}

// BitboardIterator is an efficient iterator over the contents of a bitboard.
type BitboardIterator struct {
	bitboard Bitboard
}

// Next Advances the bitboard iterator to the next state, yielding the next Square
// if there is one. returns InvalidSquare and false if there are no squares
// remaining in this bitboard.
func (bi *BitboardIterator) Next() (Square, bool) {
	if bi.bitboard == 0 {
		return InvalidSquare, false
	}

	next := bits.TrailingZeros64(uint64(bi.bitboard))
	bi.bitboard &= bi.bitboard - 1
	return Square(next), true
}
