package bridge

import (
	"os/exec"
	"strings"

	"github.com/carapace-sh/carapace"
)

// ActionBpaf bridges https://github.com/pacak/bpaf
//
// Note: bpaf's completion protocol only exposes preferred (long) flag names and
// evaluates help/version after completion output is generated, so shorthands and
// `--help`/`--version` are not completed (parity with bpaf's own completions).
//
//	var rootCmd = &cobra.Command{
//		Use:                "demo",
//		Short:              "demo bpaf application",
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
//			bridge.ActionBpaf("demo"),
//		)
//	}
func ActionBpaf(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if _, err := exec.LookPath(command[0]); err != nil {
				return carapace.ActionMessage(err.Error())
			}

			args := append([]string{"--bpaf-complete-rev=9"}, command[1:]...)
			args = append(args, c.Args...)
			args = append(args, c.Value) // an empty value is the sentinel for "completing after a space"

			return carapace.ActionExecCommand(command[0], args...)(func(output []byte) carapace.Action {
				lines := strings.Split(string(output), "\n")

				vals := make([]string, 0)
				for i := len(lines) - 1; i >= 0; i-- { // fish output is in reverse order
					if lines[i] == "" {
						continue
					}
					value, description, _ := strings.Cut(lines[i], "\t")
					vals = append(vals, value, description)
				}

				if len(vals) == 0 {
					return carapace.ActionFiles()
				}

				a := carapace.ActionValuesDescribed(vals...)
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
