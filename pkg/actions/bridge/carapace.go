package bridge

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
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
		args = append(args, c.Value)
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
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			args := []string{command[0], "export", ""}
			args = append(args, command[1:]...)
			args = append(args, c.Args...)
			args = append(args, c.Value)
			return carapace.ActionExecCommand("carapace", args...)(func(output []byte) carapace.Action {
				if string(output) == "" {
					return carapace.ActionFiles()
				}
				return carapace.ActionImport(output)
			})
		})
	})
}

// TODO could use some rework - e.g. name clashing might be a little annoying
func actionCommand(command ...string) func(f func(command ...string) carapace.Action) carapace.Action {
	return func(f func(s ...string) carapace.Action) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if len(command) > 0 {
				return f(command...)
			}

			cmd := &cobra.Command{DisableFlagParsing: true}
			carapace.Gen(cmd).Standalone()
			carapace.Gen(cmd).PositionalCompletion(
				carapace.Batch(
					carapace.ActionExecutables(),
					carapace.ActionFiles(),
				).ToA(),
			)
			carapace.Gen(cmd).PositionalAnyCompletion(
				carapace.ActionCallback(func(c carapace.Context) carapace.Action {
					return f(c.Args[0]).Shift(1)
				}),
			)
			return carapace.ActionExecute(cmd)
		})
	}
}
