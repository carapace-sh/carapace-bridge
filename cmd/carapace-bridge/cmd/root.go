package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-bridge/pkg/actions/bridge"
	"github.com/carapace-sh/carapace/pkg/ps"
	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "carapace-bridge",
	Short: "completion bridge",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Example: `  carapace-bridge completion:
    bash:       source <(carapace-bridge _carapace bash)
    elvish:     eval (carapace-bridge _carapace elvish | slurp)
    fish:       carapace-bridge _carapace fish | source
    nushell:    carapace-bridge _carapace nushell
    oil:        source <(carapace-bridge _carapace oil)
    powershell: carapace-bridge _carapace powershell | Out-String | Invoke-Expression
    tcsh:       eval ` + "`" + `carapace-bridge _carapace tcsh` + "`" + `
    xonsh:      exec($(carapace-bridge _carapace xonsh))
    zsh:        source <(carapace-bridge _carapace zsh)
`,
}

func Execute(version string) error {
	rootCmd.Version = version
	return rootCmd.Execute()
}

func init() {
	carapace.Gen(rootCmd)
	rootCmd.AddGroup(&cobra.Group{ID: "bridge", Title: "Bridge Commands"})
	addSubCommand("argcomplete", "bridges https://github.com/kislyuk/argcomplete", bridge.ActionArgcomplete)
	addSubCommand("argcomplete@v1", "bridges https://github.com/kislyuk/argcomplete", bridge.ActionArgcompleteV1)
	addSubCommand("bash", "bridges completions registered in bash", bridge.ActionBash)
	addSubCommand("carapace-bin", "bridges completions registered in carapace-bin", bridge.ActionCarapaceBin)
	addSubCommand("carapace", "bridges https://github.com/carapace-sh/carapace", bridge.ActionCarapace)
	addSubCommand("clap", "bridges https://github.com/clap-rs/clap", bridge.ActionClap)
	addSubCommand("click", "bridges https://github.com/pallets/click", bridge.ActionClick)
	addSubCommand("cobra", "bridges https://github.com/spf13/cobra", bridge.ActionCobra)
	addSubCommand("complete", "bridges https://github.com/posener/complete", bridge.ActionComplete)
	addSubCommand("fish", "bridges completions registered in fish", bridge.ActionFish)
	addSubCommand("inshellisense", "bridges https://github.com/microsoft/inshellisense", bridge.ActionInshellisense)
	addSubCommand("kingpin", "bridges https://github.com/alecthomas/kingpin", bridge.ActionKingpin)
	addSubCommand("macro", "bridges macros exposed with https://github.com/carapace-sh/carapace-spec", bridge.ActionMacro)
	addSubCommand("powershell", "bridges completions registered in powershell", bridge.ActionPowershell)
	addSubCommand("urfavecli", "bridges https://github.com/urfave/cli (v2)", bridge.ActionUrfavecli)
	addSubCommand("urfavecli@v1", "bridges https://github.com/urfave/cli (v3)", bridge.ActionUrfavecliV1)
	addSubCommand("yargs", "bridges https://github.com/yargs/yargs", bridge.ActionYargs)
	addSubCommand("zsh", "bridges completions registered in zsh", bridge.ActionZsh)
}

func addSubCommand(use, short string, f func(s ...string) carapace.Action) {
	cmd := &cobra.Command{
		Use:     use + " <command>[/<shell>]",
		Short:   short,
		GroupID: "bridge",
		Example: fmt.Sprintf(`  bridge <command>:
    bash:       source <(carapace-bridge %v command/bash)
    elvish:     eval (carapace-bridge %v command/elvish | slurp)
    fish:       carapace-bridge %v command/fish | source
    nushell:    carapace-bridge %v command/nushell
    oil:        source <(carapace-bridge %v command/oil)
    powershell: carapace-bridge %v command/powershell | Out-String | Invoke-Expression
    tcsh:       eval `+"`"+`carapace-bridge %v command/tcsh`+"`"+`
    xonsh:      exec($(carapace-bridge %v command/xonsh))
    zsh:        source <(carapace-bridge %v command/zsh)
`,
			use, use, use, use, use, use, use, use, use),
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			splitted := strings.Split(args[0], "/")
			args[0] = splitted[0]
			shell := "export"
			if len(splitted) > 1 {
				shell = splitted[1]
			}

			switch len(args) {
			case 1:
				if shell == "export" {
					shell = ps.DetermineShell()
				}

				cmd := &cobra.Command{Use: splitted[0]}
				carapace.Gen(cmd)

				stdout := bytes.Buffer{}
				cmd.SetOut(&stdout)
				cmd.SetArgs([]string{"_carapace", shell})
				cmd.Execute()

				output := stdout.String()
				switch shell {
				case "xonsh":
					output = strings.Replace(output, fmt.Sprintf("'_carapace', '%v'", shell), fmt.Sprintf("'_carapace', '%v', '', '%v'", shell, use), -1) // xonsh callback
				default:
					output = strings.Replace(output, fmt.Sprintf("_carapace %v", shell), fmt.Sprintf("_carapace %v '' %v", shell, use), -1) // general callback
				}

				fmt.Fprint(rootCmd.OutOrStdout(), output)
			default:
				rootCmd.SetArgs(append([]string{"_carapace", shell, "", use}, args...))
				rootCmd.Execute()
			}
		},
	}

	// TODO remove/prevent help flag
	carapace.Gen(cmd).Standalone()

	carapace.Gen(cmd).PreRun(func(cmd *cobra.Command, args []string) {
		cmd.Use = strings.SplitN(cmd.Use, " ", 2)[0]
	})

	rootCmd.AddCommand(cmd)

	carapace.Gen(cmd).PositionalCompletion(
		carapace.ActionMultiParts("/", func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionExecutables()
			case 1:
				return carapace.ActionStyledValues(
					"bash", "#d35673",
					"bash-ble", "#c2039a",
					"elvish", "#ffd6c9",
					"export", style.Default,
					"fish", "#7ea8fc",
					"ion", "#0e5d6d",
					"nushell", "#29d866",
					"oil", "#373a36",
					"powershell", "#e8a16f",
					"tcsh", "#412f09",
					"xonsh", "#a8ffa9",
					"zsh", "#efda53",
				)
			default:
				return carapace.ActionValues()
			}
		}),
	)

	carapace.Gen(cmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			command := strings.Split(c.Args[0], "/")[0]
			c.Args = c.Args[1:]
			return f(command).Invoke(c).ToA()
		}),
	)
}
