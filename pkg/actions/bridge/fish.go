package bridge

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace"
	shlex "github.com/rsteube/carapace-shlex"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/rsteube/carapace/pkg/xdg"
)

// ActionFish bridges completions registered in fish
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
		args = append(args, c.Value)

		configPath := fmt.Sprintf("%v/carapace/bridge/fish/config.fish", configDir)
		if err := ensureExists(configPath); err != nil {
			return carapace.ActionMessage(err.Error())
		}

		snippet := fmt.Sprintf(`source "$__fish_data_dir/config.fish";source %#v;complete --do-complete="%v"`, configPath, shlex.Join(args)) // TODO needs custom escaping
		return carapace.ActionExecCommand("fish", "--no-config", "--command", snippet)(func(output []byte) carapace.Action {
			lines := strings.Split(string(output), "\n")

			vals := make([]string, 0)
			for _, line := range lines[:len(lines)-1] {
				splitted := strings.SplitN(line, "\t", 2)

				if len(splitted) > 1 {
					vals = append(vals, splitted...)
				} else {
					vals = append(vals, splitted[0], "")
				}
			}
			return carapace.ActionValuesDescribed(vals...).StyleF(func(s string, sc style.Context) string {
				if strings.HasPrefix(s, "--") && strings.Contains(s, "=") {
					s = strings.SplitN(s, "=", 2)[1] // assume optarg
				}
				return style.ForPath(s, sc)
			})
		}).NoSpace([]rune("/=@:.,")...)
	})
}
