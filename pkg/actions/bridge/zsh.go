package bridge

import (
	_ "embed"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-bridge/third_party/github.com/Valodim/zsh-capture-completion"
	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/carapace-sh/carapace/pkg/xdg"
	"github.com/carapace-sh/carapace/third_party/github.com/acarl005/stripansi"
)

// ActionZsh bridges completions registered in zsh
//
//	(uses custom `.zshrc` in `~/.config/carapace/bridge/zsh`)
func ActionZsh(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			args := []string{"--no-rcs", "-c", zsh.Script, "--"}
			args = append(args, command...)
			args = append(args, c.Args...)
			args = append(args, c.Value)

			configDir, err := xdg.UserConfigDir()
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}
			if err := ensureExists(configDir + "/carapace/bridge/zsh/.zshrc"); err != nil {
				return carapace.ActionMessage(err.Error())
			}

			c.Setenv("CARAPACE_BRIDGE_CONFIG_HOME", configDir)
			return carapace.ActionExecCommand("zsh", args...)(func(output []byte) carapace.Action {
				lines := strings.Split(string(output), "\r\n")
				vals := make([]string, 0)

				var unquoter = strings.NewReplacer(
					`\\`, `\`,
					`\&`, `&`,
					`\<`, `<`,
					`\>`, `>`,
					"\\`", "`",
					`\'`, `'`,
					`\"`, `"`,
					`\{`, `{`,
					`\}`, `}`,
					`\$`, `$`,
					`\#`, `#`,
					`\|`, `|`,
					`\?`, `?`,
					`\(`, `(`,
					`\)`, `)`,
					`\;`, `;`,
					`\ `, ` `,
					`\[`, `[`,
					`\]`, `]`,
					`\*`, `*`,
					`\~`, `~`,
				)

				for _, line := range lines[:len(lines)-1] {
					line = unquoter.Replace(line)
					if splitted := strings.SplitN(line, " -- ", 2); len(splitted) == 1 {
						vals = append(vals, stripansi.Strip(splitted[0]), "")
					} else {
						vals = append(vals, stripansi.Strip(splitted[0]), splitted[1])
					}
				}
				return carapace.ActionValuesDescribed(vals...).StyleF(style.ForPath)
			}).Invoke(c).ToA().NoSpace([]rune("/=@:.,")...)
		})
	})
}
