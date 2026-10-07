package bridge

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/carapace-sh/carapace"
)

// ActionSymfony bridges https://github.com/symfony/symfony (Console component)
//
// Note: Symfony's completion protocol only reports long option names, so
// shorthands are not completed (parity with Symfony's own completion scripts).
//
//	var rootCmd = &cobra.Command{
//		Use:                "console",
//		Short:              "demo symfony application",
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
//			bridge.ActionSymfony("console"),
//		)
//	}
func ActionSymfony(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if _, err := exec.LookPath(command[0]); err != nil {
				return carapace.ActionMessage(err.Error())
			}

			current := c.Value
			tokens := append([]string{}, command...)
			tokens = append(tokens, c.Args...)

			currentIndex := len(tokens) // cursor floats at the end
			if current != "" {
				tokens = append(tokens, current)
				currentIndex = len(tokens) - 1
			}

			args := []string{"_complete", "--no-interaction", "-s", "zsh", "-c", strconv.Itoa(currentIndex)}
			for _, token := range tokens {
				if token == "" {
					continue // empty tokens are skipped like in the official completion script
				}
				args = append(args, "-i"+token)
			}

			c.Setenv("SHELL_VERBOSITY", "0")
			return carapace.ActionExecCommand(command[0], args...)(func(output []byte) carapace.Action {
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

				a := carapace.ActionValuesDescribed(vals...)
				if strings.HasPrefix(current, "-") { // complete `--flag=<TAB>` by prefixing suggested values with the flag
					if flag, _, found := strings.Cut(current, "="); found {
						options := false
						for i := 0; i < len(vals); i += 2 {
							if strings.HasPrefix(vals[i], "-") {
								options = true
								break
							}
						}
						if !options {
							a = a.Prefix(flag + "=")
						}
					}
				}
				return a
			}).Invoke(c).ToA()
		})
	})
}
