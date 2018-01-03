package engine

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// This file provides functions for converting positions to-and-from
// FEN representation.

var FenEndOfFileError = errors.New("unexpected end-of-file when reading FEN")
var FenInvalidDigitError = errors.New("invalid digit in FEN string")
var FenSumToEightError = errors.New("rank in FEN string does not sum to 8")
var FenUnknownRuneOrEofError = errors.New("unknown rune or early EOF in FEN string")
var FenInvalidSideToMoveError = errors.New("invalid side to move in FEN string")
var FenInvalidCastleStatusError = errors.New("invalid castle status in FEN string")
var FenInvalidEnPassantError = errors.New("invalid en passant in FEN string")
var FenInvalidHalfmoveError = errors.New("invalid halfmove in FEN")
var FenInvalidFullmoveError = errors.New("invalid fullmove in FEN")

// MakeDefaultPosition creates a default chess position.
func MakeDefaultPosition() *Position {
	pos, err := MakePositionFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 2 1")
	if err != nil {
		panic(err)
	}

	return pos
}

// MakePositionFromFen parses a FEN string and produces a Position
// from it. If the string is not valid FEN, an error is returned.
func MakePositionFromFen(fen string) (*Position, error) {
	position := MakeEmptyPosition()
	runes := []rune(fen)
	index := 0

	// helper functions for parsing a list of runes
	peek := func() rune {
		if index >= len(runes) {
			return utf8.RuneError
		}

		return runes[index]
	}

	advance := func() {
		index++
	}

	eat := func(r rune) error {
		peeked := peek()
		switch peeked {
		case utf8.RuneError:
			return FenEndOfFileError
		case r:
			advance()
			return nil
		}

		return fmt.Errorf("expected character `%c`, got `%c`", r, peeked)
	}

	// helper functions for parsing particular productions of the FEN grammar
	eatBoard := func() error {
		// fen encodes the state of each rank individually as a sequence
		// of characters, followed by a slash to indicate the end of the rank.
		//
		// ranks are encoded from the highest rank to the lowest rank.
		for rank := Rank8; ; rank-- {
			// we'll assign pieces to files based on the characters we see. If
			// we see letters in the FEN string, they correspond to pieces that
			// we should place on the rank and file. If we see numbers, they
			// represent blank squares that we should skip.
			file := FileA
			for file <= FileH {
				entry := peek()
				if unicode.IsDigit(entry) {
					// entry must be a digit 1-8 instructing us to skip the
					// next 1-8 squares.
					if entry < '1' || entry > '8' {
						return FenInvalidDigitError
					}

					value := int(entry - 48)
					file += File(value)
					if file > 8 {
						return FenSumToEightError
					}

					advance()
					continue
				}

				// if it's not a digit, this character represents a piece.
				piece, err := MakePieceFromRune(entry)
				if err != nil {
					return FenUnknownRuneOrEofError
				}

				square := MakeSquare(rank, file)
				if err := position.AddPiece(square, piece); err != nil {
					panic("unexpected double-add of piece when parsing FEN")
				}

				advance()
				file++
			}

			// we're done here if this rank we just read was the first rank.
			if rank == Rank1 {
				break
			}

			// otherwise we need to eat another slash and keep going.
			if err := eat('/'); err != nil {
				return err
			}
		}

		return nil
	}

	eatSideToMove := func() error {
		if err := eat(' '); err != nil {
			return err
		}

		switch peek() {
		case 'w':
			position.sideToMove = White
		case 'b':
			position.sideToMove = Black
		default:
			return FenInvalidSideToMoveError
		}

		advance()
		return nil
	}

	eatCastleStatus := func() error {
		if err := eat(' '); err != nil {
			return err
		}

		if peek() == '-' {
			advance()
		} else {
			// k, q, K, or Q can appear in any order here and indicate
			// which color and side is able to castle from this position.
			for i := 0; i < 4; i++ {
				switch peek() {
				case 'K':
					position.castleStatus |= whiteOO
				case 'k':
					position.castleStatus |= blackOO
				case 'Q':
					position.castleStatus |= whiteOOO
				case 'q':
					position.castleStatus |= blackOOO
				case ' ':
					return nil
				default:
					return FenInvalidCastleStatusError
				}

				advance()
			}
		}

		return nil
	}

	eatEnPassant := func() error {
		if err := eat(' '); err != nil {
			return err
		}

		if peek() == '-' {
			position.enPassantSquare = InvalidSquare
			advance()
		} else {
			epFile, err := MakeFileFromRune(peek())
			if err != nil {
				return FenInvalidEnPassantError
			}

			advance()
			epRank, err := MakeRankFromRune(peek())
			if err != nil {
				return FenInvalidEnPassantError
			}

			advance()
			position.enPassantSquare = MakeSquare(epRank, epFile)
		}

		return nil
	}

	eatHalfmove := func() error {
		if err := eat(' '); err != nil {
			return err
		}

		var buf []rune
		for {
			next := peek()
			if !unicode.IsDigit(next) {
				if len(buf) == 0 {
					return FenInvalidHalfmoveError
				}

				break
			}

			buf = append(buf, next)
			advance()
		}

		parsed, err := strconv.ParseUint(string(buf), 10, 32)
		if err != nil {
			return FenInvalidHalfmoveError
		}

		position.halfmoveClock = uint32(parsed)
		return nil
	}

	eatFullmove := func() error {
		if err := eat(' '); err != nil {
			return err
		}

		var buf []rune
		for {
			if index >= len(runes) {
				break
			}

			next := peek()
			if !unicode.IsDigit(next) {
				return FenInvalidFullmoveError
			}

			buf = append(buf, next)
			advance()
		}

		parsed, err := strconv.ParseUint(string(buf), 10, 32)
		if err != nil {
			return FenInvalidFullmoveError
		}

		position.fullmoveClock = uint32(parsed)
		return nil
	}

	// Toplevel parsing begins here.
	if err := eatBoard(); err != nil {
		return nil, err
	}

	if err := eatSideToMove(); err != nil {
		return nil, err
	}

	if err := eatCastleStatus(); err != nil {
		return nil, err
	}

	if err := eatEnPassant(); err != nil {
		return nil, err
	}

	if err := eatHalfmove(); err != nil {
		return nil, err
	}

	if err := eatFullmove(); err != nil {
		return nil, err
	}

	return position, nil
}

// AsFen retuns a string representation of this position in FEN
// notation.
func (pos *Position) AsFen() string {
	buf := new(bytes.Buffer)
	for rank := Rank8; ; rank-- {
		emptySquares := 0
		for file := FileA; file <= FileH; file++ {
			square := MakeSquare(rank, file)
			piece, ok := pos.PieceAt(square)
			if ok {
				if emptySquares != 0 {
					fmt.Fprintf(buf, "%d", emptySquares)
				}

				fmt.Fprintf(buf, "%s", piece.String())
				emptySquares = 0
			} else {
				emptySquares++
			}
		}

		if emptySquares != 0 {
			fmt.Fprintf(buf, "%d", emptySquares)
		}

		if rank == Rank1 {
			break
		}

		fmt.Fprint(buf, "/")
	}

	fmt.Fprint(buf, " ")
	if pos.SideToMove() == White {
		fmt.Fprint(buf, "w")
	} else {
		fmt.Fprint(buf, "b")
	}

	fmt.Fprint(buf, " ")
	someoneCanCastle := false
	if pos.CanCastleKingside(White) {
		fmt.Fprint(buf, "K")
		someoneCanCastle = true
	}

	if pos.CanCastleQueenside(White) {
		fmt.Fprint(buf, "Q")
		someoneCanCastle = true
	}

	if pos.CanCastleKingside(Black) {
		fmt.Fprint(buf, "k")
		someoneCanCastle = true
	}

	if pos.CanCastleQueenside(Black) {
		fmt.Fprint(buf, "q")
		someoneCanCastle = true
	}

	if !someoneCanCastle {
		fmt.Fprint(buf, "-")
	}

	fmt.Fprint(buf, " ")
	if pos.HasEnPassantSquare() {
		fmt.Fprintf(buf, "%s", pos.EnPassantSquare())
	} else {
		fmt.Fprint(buf, "-")
	}

	fmt.Fprint(buf, " ")
	fmt.Fprintf(buf, "%d %d", pos.HalfmoveClock(), pos.FullmoveClock())
	return buf.String()
}
