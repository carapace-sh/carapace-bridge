package bridge

import (
	"encoding/json"
	"strings"

	"github.com/rsteube/carapace"
)

// ActionInshellisense bridges https://github.com/microsoft/inshellisense
func ActionInshellisense(command ...string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if len(command) == 0 {
			return carapace.ActionMessage("missing argument [ActionInshellisense]")
		}

		args := append(command, c.Args...)
		args = append(args, c.Value)
		input := strings.Join(args, " ") // TODO simple join for now as the lexer in inshellisense can't handle quotes and spaces anyway
		return carapace.ActionExecCommand("inshellisense", "complete", input)(func(output []byte) carapace.Action {
			var r struct {
				Suggestions []struct {
					Name        string
					Description string
				}
			}

			if err := json.Unmarshal(output, &r); err != nil {
				return carapace.ActionMessage(err.Error())
			}

			vals := make([]string, 0)
			for _, s := range r.Suggestions {
				if !strings.HasPrefix(c.Value, "-") && strings.HasPrefix(s.Name, "-") {
					continue
				}

				if strings.HasPrefix(s.Name, c.Value) ||
					(strings.HasPrefix(c.Value, "-") && strings.Contains(c.Value, "=")) {
					vals = append(vals, s.Name, s.Description)
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
