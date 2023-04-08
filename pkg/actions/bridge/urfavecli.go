package bridge

import (
	"strings"

	"github.com/rsteube/carapace"
)

// ActionUrfavecli bridges https://github.com/urfave/cli
func ActionUrfavecli(command ...string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if len(command) == 0 {
			return carapace.ActionMessage("missing argument [ActionUrfavecli]")
		}

		args := append(command[1:], c.Args...)
		args = append(args, c.CallbackValue)
		args = append(args, "--generate-bash-completion")
		return carapace.ActionExecCommand(command[0], args...)(func(output []byte) carapace.Action {
			lines := strings.Split(string(output), "\n")
			return carapace.ActionValues(lines[:len(lines)-1]...)
		}).NoSpace([]rune("/=@:.,")...)
	})
}
