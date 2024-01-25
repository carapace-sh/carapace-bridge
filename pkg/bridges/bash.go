package bridges

import (
	"os"
	"runtime"
	"strings"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/execlog"
)

func Bash() []string {
	if runtime.GOOS == "windows" {
		return []string{}
	}
	if _, err := execlog.LookPath("bash"); err != nil {
		return []string{}
	}

	return cache("bash", func() ([]string, error) {
		unique := make(map[string]bool)
		for _, location := range []string{
			"/data/data/com.termux/files/etc/bash_completion.d",                 // termux
			"/data/data/com.termux/files/usr/share/bash-completion/completions", // termux
			"/etc/bash_completion.d",                                            // linux
			"/usr/local/etc/bash_completion.d",                                  // osx
			"/usr/local/share/bash-completion/completions",                      // osx
			"/usr/share/bash-completion/completions",                            // linux
		} {
			entries, err := os.ReadDir(location)
			if err != nil {
				carapace.LOG.Println(err.Error())
				continue
			}

			for _, entry := range entries {
				if !entry.IsDir() && !strings.HasPrefix(entry.Name(), "_") {
					unique[entry.Name()] = true
				}
			}
		}

		completers := make([]string, 0)
		for c := range filter(unique, bashBuiltins) {
			completers = append(completers, c)
		}
		return completers, nil
	})
}
