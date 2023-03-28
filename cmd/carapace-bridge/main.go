package main

import "github.com/rsteube/carapace-bridge/cmd/carapace-bridge/cmd"

var version = "develop"

func main() {
	cmd.Execute(version)
}
