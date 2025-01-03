package main

import (
	"fmt"
	"strings"

	"github.com/carapace-sh/carapace-bridge/cmd/carapace-bridge/cmd"
)

var commit, date string
var version = "develop"

func main() {
	if strings.Contains(version, "SNAPSHOT") {
		version += fmt.Sprintf(" (%v) [%v]", date, commit)
	}
	cmd.Execute(version)
}
