package cmd

import (
	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-bridge/pkg/actions/bridge"
	"github.com/carapace-sh/carapace-bridge/pkg/choice"
	"github.com/spf13/cobra"
)

var chooseCmd = &cobra.Command{
	Use:   "choose variant...",
	Short: "choose variants",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch cmd.Flag("delete").Changed {
		case true:
			for _, arg := range args {
				if err := choice.Unset(arg); err != nil {
					return err
				}
			}
		default:
			for _, arg := range args {
				if err := choice.Set(choice.Parse(arg)); err != nil {
					return err
				}
			}
		}
		return nil
	},
}

func init() {
	carapace.Gen(chooseCmd).Standalone()
	chooseCmd.Flags().SetInterspersed(false)

	chooseCmd.Flags().BoolP("delete", "d", false, "delete given choice(s)")
	rootCmd.AddCommand(chooseCmd)

	carapace.Gen(chooseCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if chooseCmd.Flag("delete").Changed {
				return bridge.ActionChoices()
			}

			return carapace.ActionMultiPartsN("/", 2, func(c carapace.Context) carapace.Action {
				switch len(c.Parts) {
				case 0:
					return carapace.ActionExecutables().Suffix("/")
				default:
					// TODO highlight known bridges
					return bridge.ActionBridges().Suffix("@bridge")
				}
			})
		}),
	)
}
