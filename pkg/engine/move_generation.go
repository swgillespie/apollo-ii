package engine

func generatePawnMoves(pos *Position, moves []Move) []Move {
	addMove := func(m Move) {
		moves = append(moves, m)
	}

	color := pos.SideToMove()
	enemyPieceMap := pos.Color(color.Toggle())
	alliedPieceMap := pos.Color(color)
	allPieces := enemyPieceMap | alliedPieceMap

	var startingRank Rank
	var promoRank Rank
	var pawnDirection Direction
	var enPassantDirection Direction
	if color == White {
		startingRank = Rank2
		promoRank = Rank8
		pawnDirection = North
		enPassantDirection = South
	} else {
		startingRank = Rank7
		promoRank = Rank1
		pawnDirection = South
		enPassantDirection = North
	}

	pawns := pos.Pawns(color).Iter()
	for pawn, hasNext := pawns.Next(); hasNext; pawn, hasNext = pawns.Next() {
		// the general pawn move is that it moves one square in the pawn
		// direction
		target := pawn.Towards(pawnDirection)

		// non-capturing moves
		if target.Rank() == promoRank && !allPieces.Test(target) {
			addMove(MakePromotionMove(pawn, target, Bishop))
			addMove(MakePromotionMove(pawn, target, Knight))
			addMove(MakePromotionMove(pawn, target, Rook))
			addMove(MakePromotionMove(pawn, target, Queen))
		} else if !allPieces.Test(target) {
			addMove(MakeQuietMove(pawn, target))
		}

		// double-pawn pushes, for pawns still on their starting square
		if pawn.Rank() == startingRank {
			twoPushTarget := target.Towards(pawnDirection)
			if !allPieces.Test(target) && !allPieces.Test(twoPushTarget) {
				addMove(MakeDoublePawnPushMove(pawn, twoPushTarget))
			}
		}

		// non-ep capturing moves
		pawnAttacks := PawnAttacks(pawn, color).Iter()
		for pawnAttack, hasNextAttack := pawnAttacks.Next(); hasNextAttack; pawnAttack, hasNextAttack = pawnAttacks.Next() {
			if enemyPieceMap.Test(pawnAttack) {
				if pawnAttack.Rank() == promoRank {
					addMove(MakePromotionCaptureMove(pawn, pawnAttack, Bishop))
					addMove(MakePromotionCaptureMove(pawn, pawnAttack, Knight))
					addMove(MakePromotionCaptureMove(pawn, pawnAttack, Rook))
					addMove(MakePromotionCaptureMove(pawn, pawnAttack, Queen))
				} else {
					addMove(MakeCaptureMove(pawn, pawnAttack))
				}
			}
		}

		// en-passant
		if pos.HasEnPassantSquare() {
			epSquare := pos.EnPassantSquare()
			// would this be a normal legal attack for this pawn?
			if PawnAttacks(pawn, color).Test(epSquare) {
				// the attack square is directly behind the pawn that was pushed
				attackSquare := epSquare.Towards(enPassantDirection)
				addMove(MakeEnPassantMove(pawn, attackSquare))
			}
		}
	}

	return moves
}

func generateKnightMoves(pos *Position, moves []Move) []Move {
	addMove := func(m Move) {
		moves = append(moves, m)
	}

	color := pos.SideToMove()
	enemyPieceMap := pos.Color(color.Toggle())
	alliedPieceMap := pos.Color(color)
	knights := pos.Knights(color).Iter()
	for knight, next := knights.Next(); next; knight, next = knights.Next() {
		attacks := KnightAttacks(knight).Iter()
		for knightAttack, next := attacks.Next(); next; knightAttack, next = attacks.Next() {
			if enemyPieceMap.Test(knightAttack) {
				addMove(MakeCaptureMove(knight, knightAttack))
			} else if !alliedPieceMap.Test(knightAttack) {
				addMove(MakeQuietMove(knight, knightAttack))
			}
		}
	}

	return moves
}

func generateSlidingMoves(pos *Position,
	moves []Move,
	attackFunc func(Square, Bitboard) Bitboard,
	boardFunc func(Color) Bitboard) []Move {
	addMove := func(m Move) {
		moves = append(moves, m)
	}

	color := pos.SideToMove()
	enemyPieceMap := pos.Color(color.Toggle())
	alliedPieceMap := pos.Color(color)
	pieces := boardFunc(color).Iter()
	for piece, next := pieces.Next(); next; piece, next = pieces.Next() {
		attacks := attackFunc(piece, enemyPieceMap|alliedPieceMap).Iter()
		for attack, next := attacks.Next(); next; attack, next = pieces.Next() {
			// in theory we only need to test the end of rays
			// for occupancy...
			if enemyPieceMap.Test(attack) {
				addMove(MakeCaptureMove(piece, attack))
			} else if !alliedPieceMap.Test(attack) {
				addMove(MakeQuietMove(piece, attack))
			}
		}
	}

	return moves
}

func generateKingMoves(pos *Position, moves []Move) []Move {
	addMove := func(m Move) {
		moves = append(moves, m)
	}

	color := pos.SideToMove()
	enemyPieceMap := pos.Color(color.Toggle())
	alliedPieceMap := pos.Color(color)
	allPieces := enemyPieceMap | alliedPieceMap

	// there should only be one king but i guess it's cool to have an engine
	// that can play chess variants with multiple kings
	kings := pos.Kings(color).Iter()
	for king, next := kings.Next(); next; king, next = kings.Next() {
		attacks := KingAttacks(king).Iter()
		for attack, next := attacks.Next(); next; attack, next = attacks.Next() {
			if enemyPieceMap.Test(attack) {
				addMove(MakeCaptureMove(king, attack))
			} else if !alliedPieceMap.Test(attack) {
				addMove(MakeQuietMove(king, attack))
			}

			// pseudo-legality as as a concept breaks down a little in the
			// presence of castling.
			//
			// this move generator attempts to delegate the harder aspects
			// of move generation (e.g. detection of absolute pins) to board
			// evaluation, where we can clearly observe that any move that leads
			// directly to the capture of a king is illegal. However, we /do/
			// need to enforce castling rules here, because it is not immediately
			// obvious during board evaluation that castling rules were broken
			// in a previous move.
			//
			// therefore, there are two rules that are enforced here:
			//  1. the king can't castle out of check
			//  2. the king can't castle through check (no square that the king
			//     "slides" over can be checked)
			//
			// we can do this check efficiently using our attack bitboards.
			if !pos.IsCheck(color) {
				if pos.CanCastleKingside(color) {
					one := king.Towards(East)
					two := one.Towards(East)
					if !allPieces.Test(one) && !allPieces.Test(two) {
						if pos.SquaresAttacking(color.Toggle(), one).Empty() &&
							pos.SquaresAttacking(color.Toggle(), two).Empty() {
							addMove(MakeKingsideCastleMove(king, two))
						}
					}
				}

				if pos.CanCastleQueenside(color) {
					one := king.Towards(West)
					two := one.Towards(West)
					three := two.Towards(West)

					// three can be checked, but it can't be occupied. this is because
					// the rook needs to move "across" three, but the king does not.
					if !allPieces.Test(one) && !allPieces.Test(two) && !allPieces.Test(three) {
						if pos.SquaresAttacking(color.Toggle(), one).Empty() &&
							pos.SquaresAttacking(color.Toggle(), two).Empty() {
							addMove(MakeQueensideCastleMove(king, two))
						}
					}
				}
			}
		}
	}

	return moves
}

func generatePseudolegalMoves(pos *Position) []Move {
	moves := make([]Move, 0, options.moveGenerationBufferSize)
	moves = generatePawnMoves(pos, moves)
	moves = generateKnightMoves(pos, moves)
	moves = generateSlidingMoves(pos, moves, BishopAttacks, pos.Bishops)
	moves = generateSlidingMoves(pos, moves, RookAttacks, pos.Rooks)
	moves = generateSlidingMoves(pos, moves, QueenAttacks, pos.Queens)
	return moves
}
