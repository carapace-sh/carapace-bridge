package bridge

import (
	"slices"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-bridge/pkg/bridges"
	"github.com/carapace-sh/carapace-bridge/pkg/choices"
	"github.com/carapace-sh/carapace-bridge/pkg/env"
	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/carapace-sh/carapace/third_party/golang.org/x/sys/execabs"
)

var bridgeActions = map[string]func(command ...string) carapace.Action{
	"argcomplete":    ActionArgcomplete,
	"argcomplete_v1": ActionArgcompleteV1,
	"aws":            ActionAws,
	"bash":           ActionBash,
	"carapace":       ActionCarapace,
	"carapace-bin":   ActionCarapaceBin,
	"clap":           ActionClap,
	"click":          ActionClick,
	"cobra":          ActionCobra,
	"complete":       ActionComplete,
	"fish":           ActionFish,
	"gcloud":         ActionGcloud,
	"inshellisense":  ActionInshellisense,
	"jj":             ActionJJ,
	"kingpin":        ActionKingpin,
	"kitten":         ActionKitten,
	"powershell":     ActionPowershell,
	"urfavecli":      ActionUrfavecli,
	"urfavecli_v1":   ActionUrfavecliV1,
	"yargs":          ActionYargs,
	"zsh":            ActionZsh,
}

// TODO experimental
func Get(name string) (func(command ...string) carapace.Action, bool) {
	a, ok := bridgeActions[name]
	return a, ok
}

// ActionBridges completes available bridges.
//
//	complete
//	cobra
func ActionBridges(name string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		return carapace.Batch(
			carapace.ActionValuesDescribed(
				"argcomplete", "bridges https://github.com/kislyuk/argcomplete",
				"argcomplete_v1", "bridges https://github.com/kislyuk/argcomplete",
				"carapace", "bridges https://github.com/carapace-sh/carapace",
				"clap", "bridges https://github.com/clap-rs/clap",
				"click", "bridges https://github.com/pallets/click",
				"cobra", "bridges https://github.com/spf13/cobra",
				"complete", "bridges https://github.com/posener/complete",
				"kingpin", "bridges https://github.com/alecthomas/kingpin",
				"macro", "bridges macros exposed with https://github.com/carapace-sh/carapace-spec",
				"urfavecli", "bridges https://github.com/urfave/cli (v2)",
				"urfavecli_v1", "bridges https://github.com/urfave/cli (v3)",
				"yargs", "bridges https://github.com/yargs/yargs",
			).Style(style.Dim),
			carapace.ActionValuesDescribed(
				"aws", "bridges https://github.com/aws/aws-cli",
				"bash", "bridges completions registered in bash",
				"carapace-bin", "bridges completions registered in carapace-bin",
				"fish", "bridges completions registered in fish",
				"gcloud", "bridges https://docs.cloud.google.com/sdk/gcloud",
				"inshellisense", "bridges https://github.com/microsoft/inshellisense",
				"jj", "bridges bridges https://www.jj-vcs.dev",
				"kitten", "bridges https://github.com/kovidgoyal/kitty",
				"powershell", "bridges completions registered in powershell",
				"zsh", "bridges completions registered in zsh",
			).StyleF(func(s string, sc style.Context) string {
				executable := map[string]string{
					"aws":           "aws",
					"bash":          "bash",
					"carapace-bin":  "carapace",
					"fish":          "fish",
					"gcloud":        "gcloud",
					"jj":            "jj",
					"inshellisense": "inshellisense",
					"kitten":        "kitten",
					"powershell":    "pwsh",
					"zsh":           "zsh",
				}[s]

				if _, err := execabs.LookPath(executable); err != nil {
					return style.Carapace.KeywordNegative
				}

				switch s {
				case "bash":
					if slices.Contains(bridges.Bash(), name) {
						return style.Carapace.KeywordPositive
					}
				case "fish":
					if slices.Contains(bridges.Fish(), name) {
						return style.Carapace.KeywordPositive
					}
				case "inshellisense":
					if slices.Contains(bridges.Inshellisense(), name) {
						return style.Carapace.KeywordPositive
					}
				case "zsh":
					if slices.Contains(bridges.Zsh(), name) {
						return style.Carapace.KeywordPositive
					}
				}
				return style.Dim
			}),
		).ToA().Tag("bridges")
	})
}

// ActionBridge bridges completions defined by choices and CARAPACE_BRIDGE
func ActionBridge(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if choice, err := choices.Get(command[0]); err == nil && (choice.Group == "" || choice.Group == "bridge") {
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
