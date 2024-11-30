package cmd

import selfupdate "github.com/carapace-sh/carapace-selfupdate"

var selfupdateCmd = selfupdate.Command("carapace-sh", "carapace-bridge")

func init() {
	rootCmd.AddCommand(selfupdateCmd)
}
