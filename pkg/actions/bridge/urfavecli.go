package bridge

import (
	"strings"

	"github.com/carapace-sh/carapace"
)

// ActionUrfavecli bridges https://github.com/urfave/cli
func ActionUrfavecli(command ...string) carapace.Action {
	return actionUrfavecliV3(command...)
}

func actionUrfavecliV3(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			args := append(command[1:], c.Args...)
			args = append(args, c.Value)
			args = append(args, "--generate-shell-completion")
			return carapace.ActionExecCommand(command[0], args...)(func(output []byte) carapace.Action {
				lines := strings.Split(string(output), "\n")
				if len(lines) <= 1 {
					return carapace.ActionFiles()
				}
				return carapace.ActionValues(lines[:len(lines)-1]...).NoSpace([]rune("/=@:.,")...)
			})
		})
	})
}

// ActionUrfavecliV1 bridges https://github.com/urfave/cli (v1-v2)
func ActionUrfavecliV1(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			args := append(command[1:], c.Args...)
			args = append(args, c.Value)
			args = append(args, "--generate-bash-completion")
			return carapace.ActionExecCommand(command[0], args...)(func(output []byte) carapace.Action {
				lines := strings.Split(string(output), "\n")
				if len(lines) <= 1 {
					return carapace.ActionFiles()
				}
				return carapace.ActionValues(lines[:len(lines)-1]...).NoSpace([]rune("/=@:.,")...)
			})
		})
	})
}
