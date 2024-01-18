package env

import (
	"os"
	"strings"
)

const (
	CARAPACE_BRIDGES = "CARAPACE_BRIDGES" // order of implicit bridges
)

func Bridges() []string {
	if v, ok := os.LookupEnv(CARAPACE_BRIDGES); ok {
		return strings.Split(v, ",")
	}
	return []string{}
}
