package bridges

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/carapace-sh/carapace/pkg/execlog"
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
		for _, location := range bashCompletionLocations() {
			entries, err := os.ReadDir(location)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				if !entry.IsDir() && !strings.HasPrefix(entry.Name(), "_") {
					name := strings.TrimSuffix(entry.Name(), ".bash")
					unique[name] = true
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

// bashCompletionLocations return folders containing bash completion scripts
//
//	https://github.com/scop/bash-completion/blob/main/bash_completion (function _comp_load)
func bashCompletionLocations() []string {
	var locations []string

	// 1) From BASH_COMPLETION_USER_DIR or XDG_DATA_HOME:
	// User installed completions are looked up first.
	if userDirs, ok := os.LookupEnv("BASH_COMPLETION_USER_DIR"); ok {
		for userDir := range strings.SplitSeq(userDirs, string(os.PathListSeparator)) {
			locations = append(locations, fmt.Sprintf("%v/completions", userDir))
		}
	} else {
		if dataHome, ok := os.LookupEnv("XDG_DATA_HOME"); ok {
			locations = append(locations, fmt.Sprintf("%v/bash-completion/completions", dataHome))
		} else if home, err := os.UserHomeDir(); err == nil {
			locations = append(locations, fmt.Sprintf("%v/.local/share/bash-completion/completions", home))
		}
	}

	// 2) From the location of bash_completion (completions relative to the main script).
	// In system installations this often hits the right directory.
	if sourceDir, ok := os.LookupEnv("BASH_SOURCE"); ok {
		locations = append(locations, fmt.Sprintf("%v/completions", sourceDir))
	}

	// 3) From bin directories extracted from $PATH.
	// For each PATH entry ending in bin or sbin, look at ../share/bash-completion/completions.
	for pathDir := range strings.SplitSeq(os.Getenv("PATH"), string(os.PathListSeparator)) {
		locations = append(locations, fmt.Sprintf("%v/share/bash-completion/completions", pathDir))
	}

	// 4) From XDG_DATA_DIRS or system data dirs.
	if dataDirs, ok := os.LookupEnv("XDG_DATA_DIRS"); ok {
		for dataDir := range strings.SplitSeq(dataDirs, string(os.PathListSeparator)) {
			locations = append(locations, fmt.Sprintf("%v/bash-completion/completions", dataDir))
		}
	} else {
		locations = append(locations,
			"/data/data/com.termux/files/usr/share/bash-completion/completions", // termux
			"/usr/local/share/bash-completion/completions",                      // osx
			"/usr/share/bash-completion/completions",                            // linux
		)
	}

	// 5) Legacy/fallback locations (pre-XDG system dirs).
	locations = append(locations,
		"/data/data/com.termux/files/etc/bash_completion.d", // termux
		"/etc/bash_completion.d",                            // linux
		"/usr/local/etc/bash_completion.d",                  // osx
	)

	return locations
}
