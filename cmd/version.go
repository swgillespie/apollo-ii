package cmd

import (
	"github.com/spf13/cobra"
	"github.com/swgillespie/apollo-ii/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:  "version",
	Long: "Prints the version of this program",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("%s\n", version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
