package bridge

import (
	"maps"
	"slices"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-bridge/pkg/bridges"
	"github.com/carapace-sh/carapace-bridge/pkg/choices"
	"github.com/carapace-sh/carapace-bridge/pkg/env"
)

// TODO @ is now incompatible with variants and needs a new name
var bridgeActions = map[string]func(command ...string) carapace.Action{
	"argcomplete":   ActionArgcomplete,
	"argcompleteV1": ActionArgcompleteV1,
	"aws":           ActionAws,
	"bash":          ActionBash,
	"carapace":      ActionCarapace,
	"carapace-bin":  ActionCarapaceBin,
	"clap":          ActionClap,
	"click":         ActionClick,
	"cobra":         ActionCobra,
	"complete":      ActionComplete,
	"fish":          ActionFish,
	"gcloud":        ActionGcloud,
	"inshellisense": ActionInshellisense,
	"kingpin":       ActionKingpin,
	"kitten":        ActionKitten,
	"powershell":    ActionPowershell,
	"urfavecli":     ActionUrfavecli,
	"urfavecliV1":   ActionUrfavecliV1,
	"yargs":         ActionYargs,
	"zsh":           ActionZsh,
}

// TODO experimental
func Get(name string) (func(command ...string) carapace.Action, bool) {
	a, ok := bridgeActions[name]
	return a, ok
}

func ActionBridges() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		return carapace.ActionValues(slices.Collect(maps.Keys(bridgeActions))...).
			Tag("bridges")
	})
}

// Bridges bridges completions as defined by choices and CARAPACE_BRIDGE environment variable
func ActionBridge(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if choice, err := choices.Get(command[0]); err == nil && choice.Group == "bridge" {
				if action, ok := bridgeActions[choice.Variant]; ok {
					return action(command...)
				}
				return carapace.ActionMessage("unknown bridge/variant: %v", choice.Variant)
			}

			for _, b := range env.Bridges() {
				switch b {
				case "bash":
					if slices.Contains(bridges.Bash(), command[0]) {
						return ActionBash(command...)
					}
				case "fish":
					if slices.Contains(bridges.Fish(), command[0]) {
						return ActionFish(command...)
					}
				case "inshellisense":
					if slices.Contains(bridges.Inshellisense(), command[0]) {
						return ActionInshellisense(command...)
					}
				case "zsh":
					if slices.Contains(bridges.Zsh(), command[0]) {
						return ActionZsh(command...)
					}
				}
			}
			return carapace.ActionValues()
		})
	})
}
