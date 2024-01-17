package bridges

import (
	"os"
	"runtime"
	"strings"

	"github.com/rsteube/carapace/pkg/execlog"
)

func Bash() []string {
	if runtime.GOOS == "windows" {
		return []string{}
	}
	if _, err := execlog.LookPath("bash"); err != nil {
		return []string{}
	}

	// TODO handle different OS/locations
	// TODO completions provides by bash itself
	entries, err := os.ReadDir("/usr/share/bash-completion/completions")
	if err != nil {
		return []string{}
	}

	unique := make(map[string]bool)
	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasPrefix(entry.Name(), "_") {
			unique[entry.Name()] = true
		}
	}

	completers := make([]string, 0)
	for c := range filter(unique, bashBuiltins) {
		completers = append(completers, c)
	}
	return completers
}
