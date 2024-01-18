package bridge

import (
	"encoding/json"
	"strings"

	"github.com/rsteube/carapace"
	shlex "github.com/rsteube/carapace-shlex"
)

// ActionInshellisense bridges https://github.com/microsoft/inshellisense
func ActionInshellisense(command ...string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if len(command) == 0 {
			return carapace.ActionMessage("missing argument [ActionInshellisense]")
		}

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
				}
			}

			if err := json.Unmarshal(output, &r); err != nil {
				return carapace.ActionMessage(err.Error())
			}

			vals := make([]string, 0)
			for _, s := range r.Suggestions {
				for _, name := range s.AllNames {
					if !strings.HasPrefix(c.Value, "-") && strings.HasPrefix(name, "-") {
						continue
					}

					if strings.HasPrefix(name, c.Value) ||
						(strings.HasPrefix(c.Value, "-") && strings.Contains(c.Value, "=")) {
						vals = append(vals, name, s.Description)
					}
				}
			}
			a := carapace.ActionValuesDescribed(vals...)
			if strings.HasPrefix(c.Value, "-") && strings.Contains(c.Value, "=") {
				a = a.Prefix(strings.SplitAfterN(c.Value, "=", 2)[0])
			}
			return a
		})
	})
}
