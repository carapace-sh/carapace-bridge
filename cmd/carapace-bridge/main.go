package main

import "github.com/carapace-sh/carapace-bridge/cmd/carapace-bridge/cmd"

var version = "develop"

func main() {
	cmd.Execute(version)
}
