package bridge

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/rsteube/carapace/pkg/xdg"
)

func compline(args ...string) string {
	for index, arg := range args {
		args[index] = fmt.Sprintf("%#v", arg)
	}
	return fmt.Sprintf("%#v", strings.Join(args, " "))
}

// ActionFish bridges completions registered in fish shell
// (uses custom `config.fish` in “~/.config/carapace/bridge/fish`)
func ActionFish(command ...string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if len(command) == 0 {
			return carapace.ActionMessage("missing argument [ActionFish]")
		}

		configDir, err := xdg.UserConfigDir()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		args := append(command, c.Args...)
		args = append(args, c.CallbackValue)
		snippet := fmt.Sprintf(`complete --do-complete=%v`, compline(args...))

		c.Setenv("XDG_CONFIG_HOME", fmt.Sprintf("%v/carapace/bridge", configDir)) // use custom config.fish for bridge
		return carapace.ActionExecCommand("fish", "--command", snippet)(func(output []byte) carapace.Action {
			lines := strings.Split(string(output), "\n")
			nospace := false

			vals := make([]string, 0)
			for _, line := range lines[:len(lines)-1] {
				splitted := strings.SplitN(line, "\t", 2)

				if len(splitted) > 1 {
					vals = append(vals, splitted...)
				} else {
					vals = append(vals, splitted[0], "")
				}
				if value := splitted[0]; !nospace && len(value) > 0 && strings.ContainsAny(value[len(value)-1:], `/=@:.,-`) {
					nospace = true
				}

			}
			a := carapace.ActionValuesDescribed(vals...).StyleF(func(s string, sc style.Context) string {
				if strings.HasPrefix(s, "--") && strings.Contains(s, "=") {
					s = strings.SplitN(s, "=", 2)[1] // assume optarg
				}
				return style.ForPath(s, sc)
			})
			if nospace {
				a = a.NoSpace()
			}
			return a
		}).Invoke(c).ToA()
	}).Tag("bash bridge")
}
