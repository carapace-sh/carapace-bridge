package bridge

import (
	"github.com/carapace-sh/carapace"
)

// ActionGcloud bridges https://docs.cloud.google.com/sdk/gcloud
func ActionGcloud(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			// TODO patch user@instance and --flag=optarg as in gcloud completion script

			if c.Value == "-" {
				c.Value = "--" // seems shorthand flags aren't completed anyway so expand to longhand first
			}
			c.Setenv("CLOUDSDK_COMPONENT_MANAGER_DISABLE_UPDATE_CHECK", "1")
			return ActionArgcompleteV1("gcloud").Invoke(c).ToA()
		})
	})
}
