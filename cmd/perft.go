package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/swgillespie/apollo-ii/pkg/engine"
)

var depth int

// perftCmd represents the perft command
var perftCmd = &cobra.Command{
	Use:  "perft",
	Long: "Calculates the PERFT statistics for a given board position.",
	Run: func(cmd *cobra.Command, args []string) {
		engine.Initialize()

		start := time.Now()
		results, err := engine.Perft(args[0], depth)
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

func init() {
	perftCmd.Flags().IntVarP(&depth, "depth", "d", 3, "the ply depth to search to")
	rootCmd.AddCommand(perftCmd)
}
