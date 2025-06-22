package bridge

import (
	"slices"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-bridge/pkg/bridges"
	"github.com/carapace-sh/carapace-bridge/pkg/env"
)

var bridgeActions = map[string]func(command ...string) carapace.Action{
	"argcomplete":    ActionArgcomplete,
	"argcomplete@v1": ActionArgcompleteV1,
	"bash":           ActionBash,
	"carapace":       ActionCarapace,
	"clap":           ActionClap,
	"click":          ActionClick,
	"cobra":          ActionCobra,
	"complete":       ActionComplete,
	"fish":           ActionFish,
	"inshellisense":  ActionInshellisense,
	"kingpin":        ActionKingpin,
	"powershell":     ActionPowershell,
	"urfavecli":      ActionUrfavecli,
	"urfavecli@v1":   ActionUrfavecliV1,
	"yargs":          ActionYargs,
	"zsh":            ActionZsh,
}

// Bridges bridges completions as defined in bridges.yaml and CARAPACE_BRIDGE environment variable
func ActionBridges(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if bridge, ok := bridges.Config()[command[0]]; ok {
				if action, ok := bridgeActions[bridge]; ok {
					return action(command...)
				}
				return carapace.ActionMessage("unknown bridge: %v", bridge)
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
