package bridge

import (
	"strings"

	"github.com/carapace-sh/carapace"
)

// ActionKingpin bridges https://github.com/alecthomas/kingpin
//
//	var rootCmd = &cobra.Command{
//		Use:                "tsh",
//		Short:              "Teleport Command Line Client",
//		Run:                func(cmd *cobra.Command, args []string) {},
//		DisableFlagParsing: true,
//	}
//
//	func Execute() error {
//		return rootCmd.Execute()
//	}
//
//	func init() {
//		carapace.Gen(rootCmd).Standalone()
//
//		carapace.Gen(rootCmd).PositionalAnyCompletion(
//			bridge.ActionKingpin("tsh"),
//		)
//	}
func ActionKingpin(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			args := []string{"--completion-bash"}
			args = append(args, command[1:]...)
			args = append(args, c.Args...)
			args = append(args, c.Value)
			return carapace.ActionExecCommand(command[0], args...)(func(output []byte) carapace.Action {
				lines := strings.Split(string(output), "\n")

				if len(lines) < 2 && !strings.HasPrefix(c.Value, "-") {
					return carapace.ActionFiles()
				}

				a := carapace.ActionValues(lines...)
				for _, line := range lines {
					if len(line) > 0 && strings.ContainsAny(line[:len(line)-1], "/=@:.,") {
						a = a.NoSpace()
						break
					}
				}
				return a
			}).Invoke(c).ToA()
		})
	})
}
