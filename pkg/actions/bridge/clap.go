package bridge

import (
	"strconv"
	"strings"

	"github.com/rsteube/carapace"
)

// ActionClap bridges https://github.com/clap-rs/clap
//
//	var rootCmd = &cobra.Command{
//		Use:                "dynamic",
//		Short:              "dynamic example",
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
//			bridge.ActionClap("dynamic"),
//		)
//	}
func ActionClap(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			index := len(append(command, c.Args...))
			args := []string{"complete", "--index", strconv.Itoa(index), "--type", "9", "--no-space", "--ifs=\n", "--"}
			args = append(args, command...)
			args = append(args, c.Args...)
			args = append(args, c.Value)
			return carapace.ActionExecCommand(command[0], args...)(func(output []byte) carapace.Action {
				lines := strings.Split(string(output), "\n")

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
