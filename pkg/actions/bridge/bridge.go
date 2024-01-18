package bridge

import (
	"slices"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace-bridge/pkg/bridges"
	"github.com/rsteube/carapace-bridge/pkg/env"
)

// Bridges bridges completions as defined in CARAPACE_BRIDGE environment variable
func ActionBridges(command ...string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if len(command) == 0 {
			return carapace.ActionMessage("missing argument [ActionBridges]")
		}

		for _, b := range env.Bridges() {
			switch b {
			case "bash":
				if slices.Contains(bridges.Bash(), command[0]) { // TODO performance/caching
					return ActionBash(command...)
				}
			case "fish":
				if slices.Contains(bridges.Fish(), command[0]) { // TODO performance/caching
					ActionFish(command...)
				}
			case "inshellisense":
				if slices.Contains(bridges.Inshellisense(), command[0]) { // TODO performance/caching
					ActionInshellisense(command...)
				}
			case "zsh":
				if slices.Contains(bridges.Zsh(), command[0]) { // TODO performance/caching
					return ActionZsh(command...)
				}
			}
		}
		return carapace.ActionValues()
	})
}
