package bridges

import (
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/execlog"
)

func Zsh() []string {
	if runtime.GOOS == "windows" {
		return []string{}
	}

	if _, err := execlog.LookPath("zsh"); err != nil {
		return []string{}
	}

	return cache("zsh", func() ([]string, error) {
		out, err := execlog.Command("zsh", "--no-rcs", "-c", "printf '%s\n' $fpath").Output()
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(out), "\n")

		unique := make(map[string]bool)
		for _, line := range lines {
			entries, err := os.ReadDir(line)
			if err != nil {
				carapace.LOG.Println(err.Error())
				continue
			}
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasPrefix(entry.Name(), "_") {
					unique[strings.TrimPrefix(entry.Name(), "_")] = true
				}
			}
		}

		completers := make([]string, 0)
		for name := range filter(unique, zshBuiltins) {
			completers = append(completers, name)
		}
		sort.Strings(completers)

		return completers, nil
	})
}
