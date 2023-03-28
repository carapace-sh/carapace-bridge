package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace-bridge/pkg/actions/bridge"
	"github.com/rsteube/carapace-bridge/pkg/actions/os"
	"github.com/spf13/cobra"
)

var fishCmd = &cobra.Command{
	Use:   "fish",
	Short: "bridge fish completion",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.SetArgs(append([]string{"_carapace", "export", "", "fish"}, args...))
		rootCmd.Execute()
	},
	DisableFlagParsing: true,
}

func init() {
	carapace.Gen(fishCmd).Standalone()

	rootCmd.AddCommand(fishCmd)

	carapace.Gen(fishCmd).PositionalCompletion(
		os.ActionPathExecutables(),
	)

	carapace.Gen(fishCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			command := c.Args[0]
			c.Args = c.Args[1:]
			return bridge.ActionFish(command).Invoke(c).ToA()
		}),
	)
}
