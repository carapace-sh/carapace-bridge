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
//	https://github.com/scop/bash-completion/blob/7864377cacb7858893a475e69c9c7461c2727e02/bash_completion#L3159C5-L3194C61
func bashCompletionLocations() []string {
	locations := []string{
		// TODO fix order
		"/data/data/com.termux/files/etc/bash_completion.d", // termux
		"/etc/bash_completion.d",                            // linux
		"/usr/local/etc/bash_completion.d",                  // osx
	}

	// # Lookup order:
	// # 1) From BASH_COMPLETION_USER_DIR (e.g. ~/.local/share/bash-completion):
	// # User installed completions.
	// if [[ ${BASH_COMPLETION_USER_DIR-} ]]; then
	//     _comp_split -F : paths "$BASH_COMPLETION_USER_DIR" &&
	//         dirs+=("${paths[@]/%//completions}")
	// else
	//     dirs=("${XDG_DATA_HOME:-$HOME/.local/share}/bash-completion/completions")
	// fi
	if userDirs, ok := os.LookupEnv("BASH_COMPLETION_USER_DIR"); ok {
		for userDir := range strings.Split(userDirs, string(os.PathListSeparator)) {
			locations = append(locations, fmt.Sprintf("%v/completions", userDir))
		}
	} else {
		if dataHome, ok := os.LookupEnv("XDG_DATA_HOME"); ok {
			locations = append(locations, fmt.Sprintf("%v/.local/share/bash-completion/completions", dataHome))
		} else if home, err := os.UserHomeDir(); err == nil {
			locations = append(locations, fmt.Sprintf("%v/.local/share/bash-completion/completions", home))
		}
	}

	// # 2) From the location of bash_completion: Completions relative to the main
	// # script. This is primarily for run-in-place-from-git-clone setups, where
	// # we want to prefer in-tree completions over ones possibly coming with a
	// # system installed bash-completion. (Due to usual install layouts, this
	// # often hits the correct completions in system installations, too.)
	// if [[ $BASH_SOURCE == */* ]]; then
	//     dirs+=("${BASH_SOURCE%/*}/completions")
	// else
	//     dirs+=(./completions)
	// fi
	if sourceDir, ok := os.LookupEnv("BASH_SOURCE"); ok {
		locations = append(locations, fmt.Sprintf("%v/completions", sourceDir))
	} else {
		locations = append(locations, "./completions")
	}

	// # 3) From bin directories extracted from the specified path to the command,
	// # the real path to the command, and $PATH
	// paths=()
	// [[ $cmd == /* ]] && paths+=("${cmd%/*}")
	// _comp_realcommand "$cmd" && paths+=("${REPLY%/*}")
	// _comp_split -aF : paths "$PATH"
	// for dir in "${paths[@]%/}"; do
	//     [[ $dir == ?*/@(bin|sbin) ]] &&
	//         dirs+=("${dir%/*}/share/bash-completion/completions")
	// done
	for _, pathDir := range strings.Split(os.Getenv("$PATH"), string(os.PathListSeparator)) {
		locations = append(locations, fmt.Sprintf("%v/share/bash-completion/completions", pathDir))
	}

	// # 4) From XDG_DATA_DIRS or system dirs (e.g. /usr/share, /usr/local/share):
	// # Completions in the system data dirs.
	// _comp_split -F : paths "${XDG_DATA_DIRS:-/usr/local/share:/usr/share}" &&
	//     dirs+=("${paths[@]/%//bash-completion/completions}")
	if dataDirs, ok := os.LookupEnv("XDG_DATA_DIRS"); ok {
		for _, dataDir := range strings.Split(dataDirs, string(os.PathListSeparator)) {
			locations = append(locations, fmt.Sprintf("%v/bash-completion/completions", dataDir))
		}
	} else {
		locations = append(locations,
			"/data/data/com.termux/files/usr/share/bash-completion/completions", // termux
			"/usr/local/share/bash-completion/completions",                      // osx
			"/usr/share/bash-completion/completions",                            // linux
		)
	}

	return locations
}
