package cmd

import (
	"encoding/json"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/swgillespie/apollo-ii/pkg/engine"
	"github.com/swgillespie/apollo-ii/pkg/perft"
)

var depth int
var saveIntermediates bool

// perftCmd represents the perft command
var perftCmd = &cobra.Command{
	Use:  "perft",
	Long: "Calculates the PERFT statistics for a given board position.",
	Run: func(cmd *cobra.Command, args []string) {
		engine.Initialize()
		if saveIntermediates {
			// this mode is slightly different than the normal perft in that
			// it serializes a bunch of information about the state of the
			// game as it traverses the tree of board positions. it's primarily
			// used to debug the move generator.
			err := doIntermediatePerft(args[0], depth)
			if err != nil {
				cmd.Printf("fatal error: %s\n", err.Error())
			}

			return
		}

		start := time.Now()
		results, err := perft.Perft(args[0], depth)
		elapsed := time.Since(start)
		if err != nil {
			cmd.Printf("fatal error: %s\n", err.Error())
			return
		}

		cmd.Printf("PERFT of depth %d on position `%s`\n", depth, args[0])
		cmd.Printf("nodes:       %d\n", results.Nodes)
		cmd.Printf("captures:    %d\n", results.Captures)
		cmd.Printf("en-passants: %d\n", results.EnPassants)
		cmd.Printf("castles:     %d\n", results.Castles)
		cmd.Printf("promotions:  %d\n", results.Promotions)
		cmd.Printf("checks:      %d\n", results.Checks)
		cmd.Printf("checkmates:  %d\n", results.Checkmates)
		cmd.Printf("\ntime elapsed: %s\n", elapsed)
	},
	Args: cobra.ExactArgs(1),
}

type intermediatePosition struct {
	Fen   string
	Moves []string
}

var intermediatePositions []intermediatePosition

func doIntermediatePerft(fenStr string, depth int) error {
	pos, err := engine.MakePositionFromFen(fenStr)
	if err != nil {
		return err
	}

	intermediatePerftImpl(pos, depth)
	jsonString, err := json.MarshalIndent(intermediatePositions, "", "  ")
	if err != nil {
		return err
	}

	os.Stdout.Write(jsonString)
	os.Stdout.WriteString("\n")
	return nil
}

func intermediatePerftImpl(pos *engine.Position, depth int) {
	if depth == 0 {
		return
	}

	fen := pos.AsFen()
	intermediate := intermediatePosition{fen, make([]string, 0)}

	for _, mov := range pos.PseudolegalMoves() {
		newPos := pos.Clone()
		toMove := pos.SideToMove()
		newPos.ApplyMove(mov)
		if !newPos.IsCheck(toMove) {
			intermediate.Moves = append(intermediate.Moves, mov.UciString())
			intermediatePerftImpl(newPos, depth-1)
		}
	}

	intermediatePositions = append(intermediatePositions, intermediate)
}

func init() {
	perftCmd.Flags().IntVarP(&depth, "depth", "d", 3, "the ply depth to search to")
	perftCmd.Flags().BoolVar(&saveIntermediates, "save-intermediates", false, "write intermediate move states to standard out")
	rootCmd.AddCommand(perftCmd)
}
