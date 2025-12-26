package choice

import (
	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-bridge/pkg/choices"
)

// ActionChoices completes choices
//
//	carapace-bridge (carapace-bridge/carapace@bridge)
//	gh (gh/cobra@bridge)
func ActionChoices() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		list, err := choices.List(true)
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}
		vals := make([]string, 0)
		for _, choice := range list {
			vals = append(vals, choice.Name, choice.Format())
		}
		return carapace.ActionValuesDescribed(vals...)
	}).Tag("choices")
}
