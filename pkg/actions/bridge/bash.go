package bridge

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/rsteube/carapace/pkg/xdg"
)

//go:embed bash.sh
var bashSnippet string

// ActionBash bridges completions registered in bash
// (uses custom `.bashrc` in “~/.config/carapace/bridge/bash`)
func ActionBash(command ...string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if len(command) == 0 {
			return carapace.ActionMessage("missing argument [ActionFish]")
		}

		configDir, err := xdg.UserConfigDir()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		// replacer := strings.NewReplacer(
		// ` `, `\ `,
		// `"`, `\""`,
		// )

		args := append(command, c.Args...)
		args = append(args, c.CallbackValue)

		// for index, arg := range args {
		// args[index] = replacer.Replace(arg)
		// }

		rcfile := fmt.Sprintf("%v/carapace/bridge/bash/.bashrc", configDir)
		c.Setenv("COMP_LINE", strings.Join(args, " "))
		c.Setenv("XDG_CONFIG_HOME", fmt.Sprintf("%v/carapace/bridge", configDir))
		return carapace.ActionExecCommand("bash", "--rcfile", rcfile, "-i", "-c", bashSnippet, strings.Join(args, " "))(func(output []byte) carapace.Action {
			lines := strings.Split(string(output), "\n")
			return carapace.ActionValues(lines[:len(lines)-1]...).StyleF(style.ForPath)
		}).Invoke(c).ToA().NoSpace([]rune("/=@:.,")...) // TODO check compopt for nospace
	})
}
