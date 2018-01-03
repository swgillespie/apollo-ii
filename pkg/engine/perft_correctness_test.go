package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// PERFT numbers for a number of predetermined board positions have been
// calculated already and can be used to test the move generator.
//
// The tables below encode these known perft numbers for specific
// positions. It is expected that the move generator will conform
// to these numbers.
type expectedPerft struct {
	fen        string
	depth      int
	nodes      uint64
	captures   uint64
	enPassants uint64
	castles    uint64
	promotions uint64
	checks     uint64
	checkmates uint64
}

var perftTests = [...]expectedPerft{
	// initial position perft tests
	expectedPerft{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		1,
		20,
		0,
		0,
		0,
		0,
		0,
		0,
	},
	expectedPerft{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		2,
		400,
		0,
		0,
		0,
		0,
		0,
		0,
	},
	expectedPerft{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		3,
		8902,
		34,
		0,
		0,
		0,
		12,
		0,
	},
}

func TestPerftCorrectness(t *testing.T) {
	Initialize()

	t.Parallel()
	for _, test := range perftTests {
		t.Run(fmt.Sprintf("perft-%s-depth-%d", test.fen, test.depth), func(tt *testing.T) {
			results, err := Perft(test.fen, test.depth)
			if !assert.NoError(tt, err) {
				tt.FailNow()
			}

			assert.Equal(tt, test.nodes, results.Nodes)
			assert.Equal(tt, test.captures, results.Captures)
			assert.Equal(tt, test.enPassants, results.EnPassants)
			assert.Equal(tt, test.castles, results.Castles)
			assert.Equal(tt, test.promotions, results.Promotions)
			assert.Equal(tt, test.checks, results.Checks)
			assert.Equal(tt, test.checkmates, results.Checkmates)
		})
	}
}
