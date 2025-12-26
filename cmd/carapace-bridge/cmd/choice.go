package cmd

import (
	"fmt"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-bridge/pkg/actions/bridge"
	"github.com/carapace-sh/carapace-bridge/pkg/actions/choice"
	"github.com/carapace-sh/carapace-bridge/pkg/choices"
	"github.com/spf13/cobra"
)

var choiceCmd = &cobra.Command{
	Use:   "choice [-d] [variant]...",
	Short: "list or edit choices",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			choices, err := choices.List(true)
			if err != nil {
				return err
			}
			for _, choice := range choices {
				fmt.Println(choice.Format())
			}
			return nil
		}

		switch cmd.Flag("delete").Changed {
		case true:
			for _, arg := range args {
				if err := choices.Unset(arg); err != nil {
					return err
				}
			}
		default:
			for _, arg := range args {
				if err := choices.Set(choices.Parse(arg)); err != nil {
					return err
				}
			}
		}
		return nil
	},
}

func init() {
	carapace.Gen(choiceCmd).Standalone()
	choiceCmd.Flags().SetInterspersed(false)

	choiceCmd.Flags().BoolP("delete", "d", false, "delete given choice(s)")
	rootCmd.AddCommand(choiceCmd)

	carapace.Gen(choiceCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if choiceCmd.Flag("delete").Changed {
				return choice.ActionChoices()
			}

			return carapace.ActionMultiPartsN("/", 2, func(c carapace.Context) carapace.Action {
				switch len(c.Parts) {
				case 0:
					return carapace.ActionExecutables().Suffix("/")
				default:
					return bridge.ActionBridges(c.Parts[0]).Filter("macro").Suffix("@bridge")
				}
			})
		}),
	)
}
