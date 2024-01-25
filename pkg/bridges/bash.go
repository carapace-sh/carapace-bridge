package bridges

import (
	"os"
	"path/filepath"
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

	unique := make(map[string]bool)
	for _, location := range []string{
		"/data/data/com.termux/files/etc/bash_completion.d",                 // termux
		"/data/data/com.termux/files/usr/share/bash-completion/completions", // termux
		"/etc/bash_completion.d",                                            // linux
		"/usr/local/etc/bash_completion.d",                                  // osx
		"/usr/local/share/bash-completion/completions",                      // osx
		"/usr/share/bash-completion/completions",                            // linux
	} {
		path, err := filepath.EvalSymlinks(location)
		if err != nil {
			continue
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() && !strings.HasPrefix(entry.Name(), "_") {
				unique[entry.Name()] = true
			}
		}
		break
	}

	completers := make([]string, 0)
	for c := range filter(unique, bashBuiltins) {
		completers = append(completers, c)
	}
	return completers
}
