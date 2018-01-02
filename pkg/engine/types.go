package engine

import (
	"errors"
	"fmt"
	"strings"
)

// A Square is a single square on the game board.
type Square uint8

const (
	A1            = Square(0)
	B1            = Square(1)
	C1            = Square(2)
	D1            = Square(3)
	E1            = Square(4)
	F1            = Square(5)
	G1            = Square(6)
	H1            = Square(7)
	A2            = Square(8)
	B2            = Square(9)
	C2            = Square(10)
	D2            = Square(11)
	E2            = Square(12)
	F2            = Square(13)
	G2            = Square(14)
	H2            = Square(15)
	A3            = Square(16)
	B3            = Square(17)
	C3            = Square(18)
	D3            = Square(19)
	E3            = Square(20)
	F3            = Square(21)
	G3            = Square(22)
	H3            = Square(23)
	A4            = Square(24)
	B4            = Square(25)
	C4            = Square(26)
	D4            = Square(27)
	E4            = Square(28)
	F4            = Square(29)
	G4            = Square(30)
	H4            = Square(31)
	A5            = Square(32)
	B5            = Square(33)
	C5            = Square(34)
	D5            = Square(35)
	E5            = Square(36)
	F5            = Square(37)
	G5            = Square(38)
	H5            = Square(39)
	A6            = Square(40)
	B6            = Square(41)
	C6            = Square(42)
	D6            = Square(43)
	E6            = Square(44)
	F6            = Square(45)
	G6            = Square(46)
	H6            = Square(47)
	A7            = Square(48)
	B7            = Square(49)
	C7            = Square(50)
	D7            = Square(51)
	E7            = Square(52)
	F7            = Square(53)
	G7            = Square(54)
	H7            = Square(55)
	A8            = Square(56)
	B8            = Square(57)
	C8            = Square(58)
	D8            = Square(59)
	E8            = Square(60)
	F8            = Square(61)
	G8            = Square(62)
	H8            = Square(63)
	InvalidSquare = Square(255)
)

// MakeSquare constructs a Square from a Rank and File.
func MakeSquare(rank Rank, file File) Square {
	// the relationship between rank, file, and square:
	//    rank = square % 8
	//    file = square / 8
	//    square = rank * 8 + file
	return Square(uint8(rank)*8 + uint8(file))
}

// MakeSquareFromString creates a Square from a FEN-style string encoding
// of a square. It returns an error if the encoding is not valid.
func MakeSquareFromString(str string) (Square, error) {
	if len(str) != 2 {
		return InvalidSquare, errors.New("square strings must be 2 runes long")
	}

	runes := []rune(str)
	file, err := MakeFileFromRune(runes[0])
	if err != nil {
		return InvalidSquare, err
	}

	rank, err := MakeRankFromRune(runes[1])
	if err != nil {
		return InvalidSquare, err
	}

	return MakeSquare(rank, file), nil
}

// Rank returns the rank of this square.
func (s Square) Rank() Rank {
	return Rank(uint8(s) >> 3)
}

// File returns the file of this square.
func (s Square) File() File {
	return File(uint8(s) & 7)
}

func (s Square) Towards(dir Direction) Square {
	return Square(int64(s) + dir.AsVector())
}

func (s Square) String() string {
	return fmt.Sprintf("%s%s", s.File(), s.Rank())
}

// A Rank represents a single rank on the chessboard.
type Rank uint8

const (
	Rank1       = Rank(0)
	Rank2       = Rank(1)
	Rank3       = Rank(2)
	Rank4       = Rank(3)
	Rank5       = Rank(4)
	Rank6       = Rank(5)
	Rank7       = Rank(6)
	Rank8       = Rank(7)
	InvalidRank = Rank(255)
)

func (r Rank) String() string {
	switch r {
	case Rank1:
		return "1"
	case Rank2:
		return "2"
	case Rank3:
		return "3"
	case Rank4:
		return "4"
	case Rank5:
		return "5"
	case Rank6:
		return "6"
	case Rank7:
		return "7"
	case Rank8:
		return "8"
	}

	panic("unknown rank")
}

// MakeRankFromRune makes a Rank from a FEN-style encoding of a rank. It returns
// an error if no such rank corresponds to the given rune.
func MakeRankFromRune(r rune) (Rank, error) {
	switch r {
	case '1':
		return Rank1, nil
	case '2':
		return Rank2, nil
	case '3':
		return Rank3, nil
	case '4':
		return Rank4, nil
	case '5':
		return Rank5, nil
	case '6':
		return Rank6, nil
	case '7':
		return Rank7, nil
	case '8':
		return Rank8, nil
	}

	return InvalidRank, errors.New("invalid rune for rank")
}

// A File represents a single file on the chessboard.
type File uint8

const (
	FileA       = File(0)
	FileB       = File(1)
	FileC       = File(2)
	FileD       = File(3)
	FileE       = File(4)
	FileF       = File(5)
	FileG       = File(6)
	FileH       = File(7)
	InvalidFile = File(255)
)

func (f File) String() string {
	switch f {
	case FileA:
		return "a"
	case FileB:
		return "b"
	case FileC:
		return "c"
	case FileD:
		return "d"
	case FileE:
		return "e"
	case FileF:
		return "f"
	case FileG:
		return "g"
	case FileH:
		return "h"
	}

	panic("unknown file")
}

// MakeFileFromRune makes a File from a FEN-style encoding of the file.
// Returns the error oif the given rune does not correspond to a file.
func MakeFileFromRune(r rune) (File, error) {
	switch r {
	case 'a':
		return FileA, nil
	case 'b':
		return FileB, nil
	case 'c':
		return FileC, nil
	case 'd':
		return FileD, nil
	case 'e':
		return FileE, nil
	case 'f':
		return FileF, nil
	case 'g':
		return FileG, nil
	case 'h':
		return FileH, nil
	}

	return InvalidFile, errors.New("invalid rune for file")
}

// A Color represents a player.
type Color uint8

const (
	White = Color(0)
	Black = Color(1)
)

func (c Color) Toggle() Color {
	if c == White {
		return Black
	}

	return White
}

// A PieceKind represents a kind of piece on the chessboard.
type PieceKind uint8

const (
	Pawn   = PieceKind(0)
	Knight = PieceKind(1)
	Bishop = PieceKind(2)
	Rook   = PieceKind(3)
	Queen  = PieceKind(4)
	King   = PieceKind(5)
)

func (pk PieceKind) String() string {
	switch pk {
	case Pawn:
		return "p"
	case Rook:
		return "r"
	case Knight:
		return "n"
	case Bishop:
		return "b"
	case Queen:
		return "q"
	case King:
		return "k"
	}

	panic("unknown PieceKind")
}

// A direction represents a cardinal direction on the game board.
type Direction int8

const (
	North     = Direction(0)
	NorthEast = Direction(1)
	East      = Direction(2)
	SouthEast = Direction(3)
	South     = Direction(4)
	SouthWest = Direction(5)
	West      = Direction(6)
	NorthWest = Direction(7)
)

func (d Direction) AsVector() int64 {
	switch d {
	case North:
		return 8
	case NorthEast:
		return 9
	case East:
		return 1
	case South:
		return -8
	case SouthWest:
		return -9
	case West:
		return -1
	case NorthWest:
		return 7
	}

	panic("unknown direction")
}

// A Piece is specific piece kind that belongs to a particular player.
type Piece struct {
	kind  PieceKind
	color Color
}

// MakePiece creates a new Piece from a kind and a color.
func MakePiece(kind PieceKind, color Color) Piece {
	return Piece{kind, color}
}

// MakePieceFromRune creates a new Piece from a FEN-style encoding of the
// piece. Returns an error if the given rune does not represent a piece.
func MakePieceFromRune(r rune) (Piece, error) {
	switch r {
	case 'p':
		return MakePiece(Pawn, Black), nil
	case 'P':
		return MakePiece(Pawn, White), nil
	case 'r':
		return MakePiece(Rook, Black), nil
	case 'R':
		return MakePiece(Rook, White), nil
	case 'n':
		return MakePiece(Knight, Black), nil
	case 'N':
		return MakePiece(Knight, White), nil
	case 'b':
		return MakePiece(Bishop, Black), nil
	case 'B':
		return MakePiece(Bishop, White), nil
	case 'q':
		return MakePiece(Queen, Black), nil
	case 'Q':
		return MakePiece(Queen, White), nil
	case 'k':
		return MakePiece(King, Black), nil
	case 'K':
		return MakePiece(King, White), nil
	}

	return MakeNullPiece(), errors.New("unknown piece from rune")
}

// MakeNullPiece returns a sentinel Piece that has no specific meaning
// but is used as a return value in some cases where no Piece is valid.
func MakeNullPiece() Piece {
	return Piece{0, 0}
}

func (p Piece) String() string {
	pieceStr := p.kind.String()
	if p.color == White {
		pieceStr = strings.ToUpper(pieceStr)
	}

	return pieceStr
}
