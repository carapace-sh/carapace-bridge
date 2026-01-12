package bridge

import (
	"strings"

	"github.com/carapace-sh/carapace"
)

// ActionJJ bridges https://www.jj-vcs.dev
func ActionJJ(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			c.Setenv("COMPLETE", "fish")
			args := []string{"--"}
			args = append(args, command...)
			args = append(args, c.Args...)
			args = append(args, c.Value)
			return carapace.ActionExecCommand(command[0], args...)(func(output []byte) carapace.Action {
				lines := strings.Split(string(output), "\n")
				vals := make([]string, 0)

				flags := strings.HasPrefix(c.Value, "-")
				for _, line := range lines[:len(lines)-1] {
					value, description, _ := strings.Cut(line, "\t")
					if flags == strings.HasPrefix(value, "-") {
						vals = append(vals, value, description)
					}
				}
				return carapace.ActionValuesDescribed(vals...)
			}).Invoke(c).ToA()
		})
	})
}
