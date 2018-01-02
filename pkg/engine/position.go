package engine

import (
	"bytes"
	"errors"
	"fmt"
)

// A Position is a single game state at a point in time. It contains
// all of the information necessary to completely describe a game of
// Chess at a particular instant.
type Position struct {
	// Bitboards for every piece kind on the board. The first dimension of
	// this two-dimensional array is the color of the piece, while the second
	// dimension is the kind of the piece.
	boardsByPiece [][]Bitboard

	// Bitboards for each color. This exists purely for efficiency reasons;
	// the contents of this array can always be calculated by or'ing the
	// contents of one dimension of boardsByPiece.
	boardsByColor []Bitboard

	// The current en passant square, if an en passant move is legal from
	// this position, or InvalidSquare if no such move is legal.
	enPassantSquare Square

	// Clocks for draws by repetition.
	halfmoveClock, fullmoveClock uint32

	// The color whose turn it is to move.
	sideToMove Color

	// The castling status of the game.
	castleStatus uint8
}

// MakeEmptyPosition creates a new position representing an empty board
// with all state set to their defaults.
func MakeEmptyPosition() *Position {
	pos := &Position{
		boardsByPiece:   make([][]Bitboard, 2),
		boardsByColor:   make([]Bitboard, 2),
		enPassantSquare: InvalidSquare,
		halfmoveClock:   0,
		fullmoveClock:   0,
		sideToMove:      White,
		castleStatus:    0}

	pos.boardsByPiece[White] = make([]Bitboard, 6)
	pos.boardsByPiece[Black] = make([]Bitboard, 6)
	return pos
}

// PieceAt returns the Piece that resides at the given square, if one exists.
func (p *Position) PieceAt(square Square) (Piece, bool) {
	var color Color
	if p.boardsByColor[White].Test(square) {
		color = White
	} else if p.boardsByColor[Black].Test(square) {
		color = Black
	} else {
		return MakeNullPiece(), false
	}

	piecesBoard := p.boardsByPiece[color]
	for piece := Pawn; piece <= King; piece++ {
		if piecesBoard[piece].Test(square) {
			return MakePiece(piece, color), true
		}
	}

	// if we get here, we failed to update a bitboard somewhere.
	panic("invalid bitboard state in PieceAt!")
}

// AddPiece Adds a piece to the game board. Returns an error if a piece already
// resides at the given location.
func (p *Position) AddPiece(square Square, piece Piece) error {
	if _, hasPiece := p.PieceAt(square); hasPiece {
		return errors.New("a piece already exists at the target square")
	}

	p.boardsByColor[piece.color].Set(square)
	p.boardsByPiece[piece.color][piece.kind].Set(square)
	return nil
}

// RemovePiece removes a piece from the game board. Returns an error if
// the square is empty.
func (p *Position) RemovePiece(square Square) error {
	if piece, hasPiece := p.PieceAt(square); hasPiece {
		p.boardsByColor[piece.color].Unset(square)
		p.boardsByPiece[piece.color][piece.kind].Unset(square)
		return nil
	}

	return errors.New("the target square is empty")
}

func (p *Position) Pieces(kind PieceKind, color Color) Bitboard {
	return p.boardsByPiece[color][kind]
}

func (p *Position) Pawns(color Color) Bitboard {
	return p.Pieces(Pawn, color)
}

func (p *Position) Rooks(color Color) Bitboard {
	return p.Pieces(Rook, color)
}

func (p *Position) Knights(color Color) Bitboard {
	return p.Pieces(Knight, color)
}

func (p *Position) Bishops(color Color) Bitboard {
	return p.Pieces(Bishop, color)
}

func (p *Position) Queens(color Color) Bitboard {
	return p.Pieces(Queen, color)
}

func (p *Position) Kings(color Color) Bitboard {
	return p.Pieces(King, color)
}

func (p *Position) White() Bitboard {
	return p.Color(White)
}

func (p *Position) Black() Bitboard {
	return p.Color(Black)
}

func (p *Position) Color(color Color) Bitboard {
	return p.boardsByColor[color]
}

func (p *Position) EnPassantSquare() Square {
	return p.enPassantSquare
}

func (p *Position) HasEnPassantSquare() bool {
	return p.enPassantSquare != InvalidSquare
}

func (p *Position) HalfmoveClock() uint32 {
	return p.halfmoveClock
}

func (p *Position) FullmoveClock() uint32 {
	return p.fullmoveClock
}

func (p *Position) SideToMove() Color {
	return p.sideToMove
}

func (p *Position) CanCastleKingside(color Color) bool {
	if color == White {
		return (p.castleStatus & whiteOO) == whiteOO
	} else {
		return (p.castleStatus & blackOO) == blackOO
	}
}

func (p *Position) CanCastleQueenside(color Color) bool {
	if color == White {
		return (p.castleStatus & whiteOOO) == whiteOOO
	} else {
		return (p.castleStatus & blackOOO) == blackOOO
	}
}

// SquaresAttacking returns a bitboard of pieces of the given color
// that are currently attacking the given square. This is useful for detecting
// pins or checked squares.
func (p *Position) SquaresAttacking(color Color, square Square) Bitboard {
	occupancy := p.Color(White) | p.Color(Black)
	result := EmptyBitboard

	queens := p.Queens(color).Iter()
	for queen, next := queens.Next(); next; queen, next = queens.Next() {
		if QueenAttacks(queen, occupancy).Test(square) {
			result.Set(queen)
		}
	}

	rooks := p.Rooks(color).Iter()
	for rook, next := rooks.Next(); next; rook, next = rooks.Next() {
		if RookAttacks(rook, occupancy).Test(square) {
			result.Set(rook)
		}
	}

	bishops := p.Bishops(color).Iter()
	for bishop, next := bishops.Next(); next; bishop, next = bishops.Next() {
		if BishopAttacks(bishop, occupancy).Test(square) {
			result.Set(bishop)
		}
	}

	knights := p.Knights(color).Iter()
	for knight, next := knights.Next(); next; knight, next = knights.Next() {
		if KnightAttacks(knight).Test(square) {
			result.Set(knight)
		}
	}

	pawns := p.Pawns(color).Iter()
	for pawn, next := pawns.Next(); next; pawn, next = pawns.Next() {
		if PawnAttacks(pawn, color).Test(square) {
			result.Set(pawn)
		}

		if p.HasEnPassantSquare() {
			if p.EnPassantSquare() == square && PawnAttacks(pawn, color).Test(square) {
				result.Set(pawn)
			}
		}
	}

	kings := p.Kings(color).Iter()
	for king, next := kings.Next(); next; king, next = kings.Next() {
		if KingAttacks(king).Test(square) {
			result.Set(king)
		}
	}

	return result
}

// IsCheck returns whether or not the given color is in check.
func (p *Position) IsCheck(color Color) bool {
	kings := p.Kings(color)
	if kings.Count() == 0 {
		// if there's no king (e.g. unit test or puzzle scenario),
		// there's no check
		return false
	}

	kingIter := kings.Iter()
	for king, next := kingIter.Next(); next; king, next = kingIter.Next() {
		if p.SquaresAttacking(color.Toggle(), king) != 0 {
			return true
		}
	}

	return false
}

// IsMovePseudoLegal performs a pseudo-legality test on the given move
// and returns whether or not it can be pseudo-legally played in the current
// game.
//
// A move is pseudo-legal if it conforms to piece movement rules on the
// given board position. Pseudo-legal moves can still be illegal if they leave
// the king in check or the piece being move is absolutely pinned.
func (p *Position) IsMovePseudoLegal(mov Move) bool {
	if mov.IsNull() {
		// null moves are trivially legal
		return true
	}

	// rule 1: there must be a piece on the source square
	sourcePiece, ok := p.PieceAt(mov.Source())
	if !ok {
		return false
	}

	// rule 2: the piece being moved is owned by the player whose
	// turn it is to move
	if sourcePiece.color != p.sideToMove {
		return false
	}

	// rule 3: if there is a piece on the target square...
	//      3.1: the piece must be owned by the other player
	//      3.2: the move must be a capture
	destinationPiece, hasPiece := p.PieceAt(mov.Destination())
	if hasPiece {
		if !mov.IsCapture() {
			return false
		}

		if destinationPiece.color == p.sideToMove {
			return false
		}
	}

	// rule 4: the move must follow the movement rules for the kind of
	// piece being moved.
	// (TODO)
	return true
}

// PseudolegalMoves generates all pseudo-legal moves available from the
// given position.
func (p *Position) PseudolegalMoves() []Move {
	return generatePseudolegalMoves(p)
}

func (p *Position) ApplyMove(mov Move) {
	if options.debugChecks && !p.IsMovePseudoLegal(mov) {
		panic("ApplyMove called on a move that is not pseudo-legal")
	}

	if mov.IsNull() {
		// quick out for null moves - don't change anything but the
		// side to move
		p.sideToMove = p.sideToMove.Toggle()
		return
	}

	movingPiece := p.pieceAtOrPanic(mov.Source())
	if options.debugChecks && movingPiece.color != p.sideToMove {
		panic("moving a piece that does not belong to the moving player")
	}

	// the basic strategy here is to remove the piece from the start square
	// and add it to the target square, removing the piece at the target
	// square if this is a capture.
	p.removePieceOrPanic(mov.Source())
	if mov.IsCapture() {
		p.applyCapture(mov)
	}

	if mov.IsCastle() {
		p.applyCastle(mov)
	}

	// pieceToAdd is the piece that will be added at the target square.
	// normally it's the piece that moved, except in the case of promotion
	// when the promoted piece is added.
	pieceToAdd := movingPiece
	if mov.IsPromotion() {
		pieceToAdd = MakePiece(mov.PromotionPiece(), movingPiece.color)
	}

	p.addPieceOrPanic(mov.Destination(), pieceToAdd)
	if mov.IsDoublePawnPush() {
		// double pawn pushes set the EP-square
		var epDir Direction
		if p.sideToMove == White {
			epDir = South
		} else {
			epDir = North
		}

		p.enPassantSquare = mov.Destination().Towards(epDir)
	} else {
		p.enPassantSquare = InvalidSquare
	}

	if p.CanCastleKingside(p.sideToMove) || p.CanCastleQueenside(p.sideToMove) {
		p.updateCastleStatus(mov, movingPiece)
	}

	p.sideToMove = p.sideToMove.Toggle()
	if mov.IsCapture() || movingPiece.kind == Pawn {
		p.halfmoveClock = 0
	} else {
		// not capturing or moving a pawn counts against the fifty
		// move rule.
		p.halfmoveClock++
	}

	if p.sideToMove == White {
		// if it's white's turn to move again, a turn has ended.
		p.fullmoveClock++
	}
}

// Clone performs a deep clone of this position, returning a new Position.
func (p *Position) Clone() *Position {
	newPos := MakeEmptyPosition()
	copy(newPos.boardsByPiece[White], p.boardsByPiece[White])
	copy(newPos.boardsByPiece[Black], p.boardsByPiece[Black])
	copy(newPos.boardsByColor, p.boardsByColor)

	newPos.enPassantSquare = p.enPassantSquare
	newPos.fullmoveClock = p.fullmoveClock
	newPos.halfmoveClock = p.halfmoveClock
	newPos.sideToMove = p.sideToMove
	newPos.castleStatus = p.castleStatus
	return newPos
}

// Subroutine for handling piece capture, since some additional checks
// are required to ensure correctness when capturing rooks on their
// starting squares. We also don't want the compiler to inline this function
// since the majority of moves aren't captures.
func (p *Position) applyCapture(mov Move) {
	// en-passant is the only case when the piece being captured
	// does not lie on the same square as the move destination.
	var targetSquare Square
	if mov.IsEnPassant() {
		var direction Direction
		if p.sideToMove == White {
			direction = South
		} else {
			direction = North
		}

		epSquare := p.EnPassantSquare()
		targetSquare = epSquare.Towards(direction)
	} else {
		targetSquare = mov.Destination()
	}

	p.removePieceOrPanic(targetSquare)

	// if we are capturing a rook that has not moved from its initial
	// state (i.e. the opponent could have used it to legally castle),
	// we have to invalidate the opponent's castling rights.
	opposingSide := p.sideToMove.Toggle()
	if p.CanCastleKingside(opposingSide) {
		var startingSquare Square
		var castleFlag uint8
		if opposingSide == White {
			startingSquare = H1
			castleFlag = whiteOO
		} else {
			startingSquare = H8
			castleFlag = blackOO
		}

		if targetSquare == startingSquare {
			// if the opponent can castle kingside and we just captured
			// a piece on the kingside rook starting square, we must
			// have just captured a rook.
			//
			// we must eliminate the kingside castle.
			p.castleStatus &= ^castleFlag
		}
	}

	// same deal for queenside castles.
	if p.CanCastleQueenside(opposingSide) {
		var startingSquare Square
		var castleFlag uint8
		if opposingSide == White {
			startingSquare = A1
			castleFlag = whiteOOO
		} else {
			startingSquare = A8
			castleFlag = blackOOO
		}

		if targetSquare == startingSquare {
			p.castleStatus &= ^castleFlag
		}
	}
}

// Subroutine for handling castling, since castle moves are encoded in a
// unique way and are generally unique in chess in that they move two pieces
// instead of one.
func (p *Position) applyCastle(mov Move) {
	// castles are encoded based on the king's start and stop position.
	// notably, the rook is not at the move destination.

	// regardless of how we are castling, the rook appears
	// adjacent to the king.
	var postCastleDir Direction
	var preCastleDir Direction
	var rookSquare Square
	if mov.IsKingsideCastle() {
		postCastleDir = West
		preCastleDir = East
		rookSquare = mov.Destination().Towards(preCastleDir)
	} else {
		postCastleDir = East
		preCastleDir = West
		rookSquare = mov.Destination().Towards(preCastleDir).Towards(preCastleDir)
	}

	newRookSquare := mov.Destination().Towards(postCastleDir)
	rook := p.pieceAtOrPanic(rookSquare)
	if rook.kind != Rook {
		panic("piece at rook castle square is not a rook")
	}

	p.removePieceOrPanic(rookSquare)
	p.addPieceOrPanic(newRookSquare, rook)
}

// Subroutine for updating the castle status if any player can still castle.
func (p *Position) updateCastleStatus(mov Move, movingPiece Piece) {
	switch movingPiece.kind {
	case King:
		// if it's the king that's moving, we can't castle in
		// either direction anymore.
		p.clearCastleStatus()
	case Rook:
		var kingsideRook, queensideRook Square
		if p.sideToMove == White {
			kingsideRook = H1
			queensideRook = A1
		} else {
			kingsideRook = H8
			queensideRook = A8
		}

		if p.CanCastleQueenside(p.sideToMove) && mov.Source() == queensideRook {
			p.clearQueensideCastle()
		}

		if p.CanCastleKingside(p.sideToMove) && mov.Source() == kingsideRook {
			p.clearKingsideCastle()
		}
	default:
		// other moves don't influence castle status.
	}
}

func (p *Position) addPieceOrPanic(square Square, piece Piece) {
	if err := p.AddPiece(square, piece); err != nil {
		panic(fmt.Sprintf("failed to add piece to board: %s\n", err.Error()))
	}
}

func (p *Position) removePieceOrPanic(square Square) {
	if err := p.RemovePiece(square); err != nil {
		panic(fmt.Sprintf("failed to remove piece from board: %s\n", err.Error()))
	}
}

func (p *Position) pieceAtOrPanic(square Square) Piece {
	if piece, ok := p.PieceAt(square); ok {
		return piece
	}

	panic("failed to retrieve piece on board")
}

func (p *Position) clearCastleStatus() {
	var mask uint8
	if p.sideToMove == White {
		mask = whiteCastleMask
	} else {
		mask = blackCastleMask
	}

	p.castleStatus &= ^mask
}

func (p *Position) clearKingsideCastle() {
	var mask uint8
	if p.sideToMove == White {
		mask = whiteOO
	} else {
		mask = blackOO
	}

	p.castleStatus &= ^mask
}

func (p *Position) clearQueensideCastle() {
	var mask uint8
	if p.sideToMove == White {
		mask = whiteOOO
	} else {
		mask = blackOOO
	}

	p.castleStatus &= ^mask
}

func (p *Position) String() string {
	buf := new(bytes.Buffer)
	for rank := Rank8; ; rank-- {
		for file := FileA; file <= FileH; file++ {
			square := MakeSquare(rank, file)
			piece, hasPiece := p.PieceAt(square)
			if hasPiece {
				fmt.Fprintf(buf, " %s ", piece.String())
			} else {
				fmt.Fprint(buf, " . ")
			}
		}

		fmt.Fprintf(buf, "| %s\n", Rank(rank).String())
		if rank == Rank1 {
			break
		}
	}

	for i := FileA; i <= FileH; i++ {
		fmt.Fprint(buf, "---")
	}

	fmt.Fprintln(buf)
	for file := FileA; file <= FileH; file++ {
		fmt.Fprintf(buf, " %s ", File(file).String())
	}

	fmt.Fprintln(buf)
	return buf.String()
}

// Flags for the castleStatus flag of a Position, indicating what castling
// moves are legal from this Position.
const (
	castleNone      = 0x0
	whiteOO         = 0x1
	whiteOOO        = 0x2
	whiteCastleMask = 0x3
	blackOO         = 0x4
	blackOOO        = 0x8
	blackCastleMask = 0xC
)
