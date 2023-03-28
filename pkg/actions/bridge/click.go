package bridge

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/rsteube/carapace"
)

// ActionClick bridges https://github.com/pallets/click
//
//	var rootCmd = &cobra.Command{
//		Use:                "watson",
//		Short:              "Watson is a tool aimed at helping you monitoring your time",
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
//			bridge.ActionClick("watson"),
//		)
//	}
func ActionClick(command ...string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if len(command) == 0 {
			return carapace.ActionMessage("missing argument [ActionClick]")
		}

		if _, err := exec.LookPath(command[0]); err != nil {
			return carapace.ActionMessage(err.Error())
		}

		args := append(command[1:], c.Args...)
		current := c.CallbackValue

		compLine := command[0] + " " + strings.Join(append(args, current), " ") // TODO escape/quote special characters
		c.Setenv(fmt.Sprintf("_%v_COMPLETE", strings.ToUpper(command[0])), "zsh_complete")
		c.Setenv("COMP_WORDS", compLine)
		c.Setenv("COMP_CWORD", strconv.Itoa(len(args)+1))
		return carapace.ActionExecCommand(command[0])(func(output []byte) carapace.Action {
			lines := strings.Split(string(output), "\n")

			vals := make([]string, 0)
			for i := 0; i+3 < len(lines); i += 3 {
				switch lines[i] { // type
				case "dir":
					return carapace.ActionDirectories()
				case "file":
					return carapace.ActionFiles()
				case "plain":
					value := lines[i+1]
					description := lines[i+2]
					if description == "_" {
						description = ""
					}
					vals = append(vals, value, description)
				}
			}
			return carapace.ActionValuesDescribed(vals...)
		}).Invoke(c).ToA()
	})
}
