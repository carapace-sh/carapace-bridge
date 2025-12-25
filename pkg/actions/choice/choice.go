package choice

import (
	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-bridge/pkg/choice"
)

func ActionChoices() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		choices, err := choice.List(true)
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}
		vals := make([]string, 0)
		for _, choice := range choices {
			vals = append(vals, choice.Name, choice.Format())
		}
		return carapace.ActionValues(vals...)
	})
}
