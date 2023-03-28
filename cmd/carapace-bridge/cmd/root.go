package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "carapace-bridge",
	Short: "",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute(version string) error {
	rootCmd.Version = version
	return rootCmd.Execute()
}
func init() {
	carapace.Gen(rootCmd).Standalone()
}
