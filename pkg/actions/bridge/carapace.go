package bridge

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
)

// ActionCarapace bridges https://github.com/carapace-sh/carapace
func ActionCarapace(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
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
	})
}

// ActionCarapaceBin bridges completions registered in carapace-bin
func ActionCarapaceBin(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			cmd := "carapace"
			if executable, err := os.Executable(); err == nil && filepath.Base(executable) == "carapace" {
				cmd = executable // workaround for sandbox tests: directly call executable which was built with "go run"
			}

			args := []string{command[0], "export", ""}
			args = append(args, command[1:]...)
			args = append(args, c.Args...)
			args = append(args, c.Value)
			return carapace.ActionExecCommand(cmd, args...)(func(output []byte) carapace.Action {
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

// ActionCarapace bridges macros exposed with https://github.com/rsteube/carapace-bin
func ActionMacro(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			args := []string{"_carapace", "macro"}
			args = append(args, command[1:]...)
			args = append(args, c.Args...)
			args = append(args, c.Value)

			switch len(append(command, c.Args...)) {
			case 1:
				return carapace.ActionExecCommand(command[0], "_carapace", "macro")(func(output []byte) carapace.Action {
					lines := strings.Split(string(output), "\n")
					return carapace.ActionValues(lines[:len(lines)-1]...).MultiParts(".")
				})
			default:
				return carapace.ActionExecCommand(command[0], args...)(func(output []byte) carapace.Action {
					return carapace.ActionImport(output)
				}).Shift(1)
			}
		})
	})
}
