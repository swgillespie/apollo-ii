package engine

import (
	"math/bits"
)

// This file encapsules the logic of attack move generation.
// For the purposes of speed, attack moves are precomputed and stored in a
// table which is then used by the move generator when generating moves or
// by the position evaluator when determining whether or not a king is in check.
//
// This module currently implements the "classic method" of move generation,
// which precomputes sliding rays of attack for every sliding piece on the
// board and every direction. Movesets for queens, rooks, and bishops can be
// constructed by taking the union of move rays in legal directions.
// Movesets for kings, pawns, and knights do not need to consider blocking
// pieces.
//
// All of the sliding functions in this module consider the first blocking
// piece along a ray to be a legal move, which it is if the first blocking
// piece is an enemy piece. It is the responsibility of callers of this
// function to determine whether or not the blocking piece is an enemy piece.
var rayTable [][]Bitboard
var pawnTable [][]Bitboard
var knightTable []Bitboard
var kingTable []Bitboard

// a ray is "positive" if the ray vector is positive, otherwise a ray is
// "negative". if a ray is negative, we need to use leading zeros intead of
// trailing zeros in order to find the blocking piece.
func positiveRayAttacks(square Square, occupancy Bitboard, direction Direction) Bitboard {
	attacks := rayTable[square][direction]
	blocker := attacks & occupancy
	if blocker == 0 {
		// no blockers - just return the attack board
		return attacks
	}

	blockingSquare := Square(bits.TrailingZeros64(uint64(blocker)))
	blockingRay := rayTable[blockingSquare][direction]
	return attacks ^ blockingRay
}

func negativeRayAttacks(square Square, occupancy Bitboard, direction Direction) Bitboard {
	attacks := rayTable[square][direction]
	blocker := attacks & occupancy
	if blocker == 0 {
		// no blockers - just return the attack board
		return attacks
	}

	blockingSquare := 63 - bits.LeadingZeros64(uint64(blocker))
	blockingRay := rayTable[blockingSquare][direction]
	return attacks ^ blockingRay
}

func diagonalAttacks(square Square, occupancy Bitboard) Bitboard {
	return positiveRayAttacks(square, occupancy, NorthWest) | negativeRayAttacks(square, occupancy, SouthEast)
}

func antidiagonalAttacks(square Square, occupancy Bitboard) Bitboard {
	return positiveRayAttacks(square, occupancy, NorthEast) | negativeRayAttacks(square, occupancy, SouthWest)
}

func fileAttacks(square Square, occupancy Bitboard) Bitboard {
	return positiveRayAttacks(square, occupancy, North) | negativeRayAttacks(square, occupancy, South)
}

func rankAttacks(square Square, occupancy Bitboard) Bitboard {
	return positiveRayAttacks(square, occupancy, East) | negativeRayAttacks(square, occupancy, West)
}

// BishopAttacks Returns the bitboard of legal bishop moves for a piece at the given square
// and with the given board occupancy.
func BishopAttacks(square Square, occupancy Bitboard) Bitboard {
	return diagonalAttacks(square, occupancy) | antidiagonalAttacks(square, occupancy)
}

// RookAttacks Returns the bitboard of legal rook moves for a piece at the given square
// and with the given board occupancy.
func RookAttacks(square Square, occupancy Bitboard) Bitboard {
	return rankAttacks(square, occupancy) | fileAttacks(square, occupancy)
}

// QueenAttacks Returns the bitboard of legal queen moves for a piece at the given square
// and with the given board occupancy.
func QueenAttacks(square Square, occupancy Bitboard) Bitboard {
	return BishopAttacks(square, occupancy) | RookAttacks(square, occupancy)
}

// KnightAttacks Returns the bitboard of legal knight moves for a piece at the given square.
func KnightAttacks(square Square) Bitboard {
	return knightTable[square]
}

// PawnAttacks Returns the bitboard of legal pawn moves for a pawn at the given square
// and with the given color.
func PawnAttacks(square Square, color Color) Bitboard {
	return pawnTable[square][color]
}

// KingAttacks Returns the bitboard of legal king moves for the given square.
func KingAttacks(square Square) Bitboard {
	return kingTable[square]
}

func populateDirection(square Square, direction Direction, edge Bitboard) {
	if edge.Test(square) {
		// nothing to do here, there are no legal moves on this ray
		// from this square.
		return
	}

	// starting at the given square, cast a ray in the given direction
	// and add all bits to the ray mask.
	entry := &rayTable[square][direction]
	for {
		cursor := Square(int64(square) + direction.AsVector())
		entry.Set(cursor)

		// did we reach the end of the board? if so, stop.
		if edge.Test(cursor) {
			break
		}
	}
}

// Initializes all of the global precomputed state required for efficient
// run-time lookups of sliding moves.
func initializeRays() {
	// the idea here is to generate rays in every direction for every square
	// on the board, to be used by the above methods.

	rank8 := FullBitboard.Rank(Rank8)
	rank1 := FullBitboard.Rank(Rank1)
	filea := FullBitboard.File(FileA)
	fileh := FullBitboard.File(FileH)
	for sq := A1; sq <= H8; sq++ {
		populateDirection(sq, North, rank8)
		populateDirection(sq, NorthEast, rank8|fileh)
		populateDirection(sq, East, fileh)
		populateDirection(sq, SouthEast, rank1|fileh)
		populateDirection(sq, South, rank1)
		populateDirection(sq, SouthWest, rank1|filea)
		populateDirection(sq, West, filea)
		populateDirection(sq, NorthWest, rank8|filea)
	}
}

func initializeKings() {
	filea := FullBitboard.File(FileA)
	fileh := FullBitboard.File(FileH)
	rank1 := FullBitboard.Rank(Rank1)
	rank8 := FullBitboard.Rank(Rank8)
	for sq := A1; sq <= H8; sq++ {
		board := EmptyBitboard
		boardRef := &board

		if !rank8.Test(sq) {
			boardRef.Set(sq.Towards(North))
			if !filea.Test(sq) {
				boardRef.Set(sq.Towards(NorthWest))
			}

			if !fileh.Test(sq) {
				boardRef.Set(sq.Towards(NorthEast))
			}
		}

		if !rank1.Test(sq) {
			boardRef.Set(sq.Towards(South))
			if !filea.Test(sq) {
				boardRef.Set(sq.Towards(SouthWest))
			}

			if !fileh.Test(sq) {
				boardRef.Set(sq.Towards(SouthEast))
			}
		}

		if !filea.Test(sq) {
			boardRef.Set(sq.Towards(West))
		}

		if !fileh.Test(sq) {
			boardRef.Set(sq.Towards(East))
		}

		kingTable[sq] = board
	}
}

// InitializeAttackTables initializes all precomputed state about attack moves.
func InitializeAttackTables() {
	initializeRays()
}
