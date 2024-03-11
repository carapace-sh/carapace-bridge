package bridge

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/carapace-sh/carapace"
	shlex "github.com/carapace-sh/carapace-shlex"
	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/carapace-sh/carapace/pkg/xdg"
)

//go:embed bash.sh
var bashSnippet string

// ActionBash bridges completions registered in bash
// (uses custom `.bashrc` in “~/.config/carapace/bridge/bash`)
func ActionBash(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			configDir, err := xdg.UserConfigDir()
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}

			args := append(command, c.Args...)
			args = append(args, c.Value)

			configPath := fmt.Sprintf("%v/carapace/bridge/bash/.bashrc", configDir)
			if err := ensureExists(configPath); err != nil {
				return carapace.ActionMessage(err.Error())
			}

			joined := shlex.Join(args)
			if c.Value == "" {
				joined = strings.TrimSuffix(joined, `""`)
			}
			c.Setenv("COMP_LINE", joined)

			file, err := os.CreateTemp(os.TempDir(), "carapace-bridge_bash_*")
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}
			defer os.Remove(file.Name())

			os.WriteFile(file.Name(), []byte(bashSnippet), os.ModePerm)

			return carapace.ActionExecCommand("bash", "--rcfile", configPath, "-i", file.Name())(func(output []byte) carapace.Action {
				lines := strings.Split(string(output), "\n")

				vals := make([]string, 0)
				for _, line := range lines[:len(lines)-1] {
					if splitted := strings.SplitN(line, "(", 2); len(splitted) == 2 {
						// assume results contain descriptions in the format `value (description)` (spf13/cobra, carapace-sh/carapace)
						vals = append(vals,
							strings.TrimSpace(splitted[0]),
							strings.TrimSpace(strings.TrimSuffix(splitted[1], ")")),
						)
					} else {
						vals = append(vals, strings.TrimSpace(line), "")
					}
				}
				switch len(vals) {
				case 0:
					return carapace.ActionFiles()
				default:
					return carapace.ActionValuesDescribed(vals...).StyleF(style.ForPath)
				}
			}).Invoke(c).ToA().NoSpace([]rune("/=@:.,")...) // TODO check compopt for nospace
		})
	})
}
