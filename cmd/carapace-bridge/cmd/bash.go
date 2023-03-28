package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var bashCmd = &cobra.Command{
	Use:   "bash",
	Short: "",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(bashCmd).Standalone()

	rootCmd.AddCommand(bashCmd)
}
