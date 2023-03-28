package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace-bridge/pkg/actions/bridge"
	"github.com/rsteube/carapace-bridge/pkg/actions/os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "carapace-bridge",
	Short: "bridge completions",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute(version string) error {
	rootCmd.Version = version
	return rootCmd.Execute()
}
func init() {
	carapace.Gen(rootCmd).Standalone()
	addSubCommand("argcomplete", "bridges https://github.com/kislyuk/argcomplete", bridge.ActionArgcomplete)
	addSubCommand("carapace-bin", "bridges completions registered in carapace-bin", bridge.ActionCarapaceBin)
	addSubCommand("carapace", "bridges https://github.com/rsteube/carapace", bridge.ActionCarapace)
	addSubCommand("click", "bridges https://github.com/pallets/click", bridge.ActionClick)
	addSubCommand("click", "bridges https://github.com/pallets/click", bridge.ActionClick)
	addSubCommand("cobra", "bridges https://github.com/spf13/cobra", bridge.ActionCobra)
	addSubCommand("complete", "bridges https://github.com/spf13/cobra", bridge.ActionComplete)
	addSubCommand("fish", "bridges completions registered in fish shell", bridge.ActionFish)
	addSubCommand("yargs", "bridges https://github.com/yargs/yargs", bridge.ActionYargs)
}

func addSubCommand(use, short string, f func(s ...string) carapace.Action) {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.SetArgs(append([]string{"_carapace", "export", "", use}, args...))
			rootCmd.Execute()
		},
		DisableFlagParsing: true,
	}

	carapace.Gen(cmd).Standalone()

	rootCmd.AddCommand(cmd)

	carapace.Gen(cmd).PositionalCompletion(
		os.ActionPathExecutables(),
	)

	carapace.Gen(cmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			command := c.Args[0]
			c.Args = c.Args[1:]
			return f(command).Invoke(c).ToA()
		}),
	)
}
