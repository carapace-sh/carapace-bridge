package bridges

import (
	"os"
	"runtime"
	"strings"

	"github.com/rsteube/carapace/pkg/execlog"
)

func Fish() []string {
	if runtime.GOOS == "windows" {
		return []string{}
	}
	if _, err := execlog.LookPath("fish"); err != nil {
		return []string{}
	}

	return cache("fish", func() ([]string, error) {
		unique := make(map[string]bool)
		for _, location := range []string{
			"/data/data/com.termux/files/usr/share/fish/completions",          // termux
			"/data/data/com.termux/files/usr/share/fish/vendor_completions.d", // termux
			"/usr/local/share/fish/completions",                               // osx
			"/usr/local/share/fish/vendor_completions.d",                      // osx
			"/usr/share/fish/completions",                                     // linux
			"/usr/share/fish/vendor_completions.d",                            // linux
		} {
			entries, err := os.ReadDir(location)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".fish") {
					unique[strings.TrimSuffix(entry.Name(), ".fish")] = true
				}
			}
		}

		completers := make([]string, 0)
		for c := range filter(unique, fishBuiltins) {
			completers = append(completers, c)
		}
		return completers, nil
	})

}
