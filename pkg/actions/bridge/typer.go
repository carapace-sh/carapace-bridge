package bridge

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/carapace-sh/carapace"
	shlex "github.com/carapace-sh/carapace-shlex"
)

// ActionTyper bridges https://github.com/fastapi/typer
//
//	var rootCmd = &cobra.Command{
//		Use:                "clab-connector",
//		Short:              "Connector for Containerlab",
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
//			bridge.ActionTyper("clab-connector"),
//		)
//	}
func ActionTyper(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if _, err := exec.LookPath(command[0]); err != nil {
				return carapace.ActionMessage(err.Error())
			}

			args := append(command[1:], c.Args...)
			current := c.Value

			// typer reads the full command line from _TYPER_COMPLETE_ARGS
			// (shlex-style split) and treats the last word as incomplete
			// unless it ends with a space (see typer/_completion_classes.py).
			compLine := command[0] + " " + shlex.Join(args)
			if current == "" {
				compLine += " "
			} else {
				compLine += " " + shlex.Join([]string{current})
			}

			autocompleteVar := fmt.Sprintf("_%v_COMPLETE", strings.ToUpper(strings.NewReplacer("-", "_", ".", "_").Replace(command[0])))
			c.Setenv(autocompleteVar, "complete_fish")
			c.Setenv("_TYPER_COMPLETE_FISH_ACTION", "get-args")
			c.Setenv("_TYPER_COMPLETE_ARGS", compLine)

			// typer outputs newline-separated tab-completed `value\tdescription`
			// pairs and an empty string when there are no completions (file
			// completion fallback in the fish script).
			return carapace.ActionExecCommand(command[0])(func(output []byte) carapace.Action {
				lines := strings.Split(string(output), "\n")
				vals := make([]string, 0)
				for _, line := range lines {
					if line == "" {
						continue
					}
					value, description, _ := strings.Cut(line, "\t")
					vals = append(vals, value, description)
				}

				if len(vals) == 0 {
					return carapace.ActionFiles()
				}
				return carapace.ActionValuesDescribed(vals...)
			}).Invoke(c).ToA()
		})
	})
}
