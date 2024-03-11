package bridge

import (
	"encoding/json"
	"strings"

	"github.com/carapace-sh/carapace"
	shlex "github.com/carapace-sh/carapace-shlex"
	"github.com/carapace-sh/carapace/pkg/style"
)

// ActionInshellisense bridges https://github.com/microsoft/inshellisense
func ActionInshellisense(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			args := append(command, c.Args...)
			args = append(args, c.Value)

			input := shlex.Join(args)

			if strings.HasSuffix(input, `""`) {
				// TODO temporary fix as inshellisense can't handle quotes yet (won't work for those within)
				input = input[:len(input)-2] + " "
			}

			return carapace.ActionExecCommand("inshellisense", "complete", input)(func(output []byte) carapace.Action {
				var r struct {
					Suggestions []struct {
						Name        string
						AllNames    []string `json:"allNames"`
						Description string
						Icon        string
					}
				}

				if err := json.Unmarshal(output, &r); err != nil {
					return carapace.ActionMessage(err.Error())
				}

				vals := make([]string, 0)
				flags := make([]string, 0)
				files := make([]string, 0)
				for _, s := range r.Suggestions {
					for _, name := range s.AllNames {
						if !strings.HasPrefix(c.Value, "-") && strings.HasPrefix(name, "-") {
							continue
						}

						if strings.HasPrefix(name, c.Value) ||
							(strings.HasPrefix(c.Value, "-") && strings.Contains(c.Value, "=")) {
							switch s.Icon {
							case "📀":
								files = append(files, name, s.Description)
							case "🔗":
								flags = append(flags, name, s.Description)
							default:
								vals = append(vals, name, s.Description)
							}
						}
					}
				}
				a := carapace.Batch(
					carapace.ActionValuesDescribed(vals...),
					carapace.ActionValuesDescribed(flags...).Tag("flags"),
					carapace.ActionValuesDescribed(files...).StyleF(style.ForPathExt).Tag("files"),
				).ToA()
				if strings.HasPrefix(c.Value, "-") && strings.Contains(c.Value, "=") {
					a = a.Prefix(strings.SplitAfterN(c.Value, "=", 2)[0])
				}
				return a
			})
		})
	})
}
