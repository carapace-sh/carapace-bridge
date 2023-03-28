package bridge

import (
	"github.com/rsteube/carapace"
)

// ActionCarapace bridges https://github.com/rsteube/carapace
func ActionCarapace(command ...string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if len(command) == 0 {
			return carapace.ActionMessage("missing argument [ActionCarapace]")
		}

		args := []string{"_carapace", "export", ""}
		args = append(args, command[1:]...)
		args = append(args, c.Args...)
		args = append(args, c.CallbackValue)
		return carapace.ActionExecCommand(command[0], args...)(func(output []byte) carapace.Action {
			if string(output) == "" {
				return carapace.ActionValues()
			}
			return carapace.ActionImport(output)
		})
	})
}

// ActionCarapaceBin bridges completions registered in carapace-bin
func ActionCarapaceBin(command ...string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if len(command) == 0 {
			return carapace.ActionMessage("missing argument [ActionCarapaceBin]")
		}

		args := []string{command[0], "export", ""}
		args = append(args, command[1:]...)
		args = append(args, c.Args...)
		args = append(args, c.CallbackValue)
		return carapace.ActionExecCommand("carapace", args...)(func(output []byte) carapace.Action {
			if string(output) == "" {
				return carapace.ActionFiles()
			}
			return carapace.ActionImport(output)
		})
	})
}
