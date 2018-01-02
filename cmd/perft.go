package cmd

import (
	"github.com/spf13/cobra"
	"github.com/swgillespie/apollo-ii/pkg/engine/perft"
)

// perftCmd represents the perft command
var perftCmd = &cobra.Command{
	Use:  "perft",
	Long: "Calculates the PERFT statistics for a given board position.",
	Run: func(cmd *cobra.Command, args []string) {
		doPerft(cmd, args[0])
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(perftCmd)
}

func doPerft(cmd *cobra.Command, fen string) {
	results, err := perft.Perft(fen, 3)
	if err != nil {
		cmd.Printf("fatal error: %s\n", err.Error())
		return
	}

	cmd.Printf("nodes: %d\n", results.Nodes)
	//cmd.Printf("captures: %d\n", results.Captures)
}
