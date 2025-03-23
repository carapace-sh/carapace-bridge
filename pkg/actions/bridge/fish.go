package bridge

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/carapace-sh/carapace"
	shlex "github.com/carapace-sh/carapace-shlex"
	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/carapace-sh/carapace/pkg/xdg"
)

// ActionFish bridges completions registered in fish
// (uses custom `config.fish` in “~/.config/carapace/bridge/fish`)
func ActionFish(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			configDir, err := xdg.UserConfigDir()
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}

			args := append(command, c.Args...)
			args = append(args, c.Value)

			configPath := filepath.Join(configDir, "carapace/bridge/fish/config.fish")
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
	})
}
