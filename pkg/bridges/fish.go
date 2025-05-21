package bridges

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/carapace-sh/carapace/pkg/execlog"
	"github.com/carapace-sh/carapace/pkg/xdg"
)

func Fish() []string {
	if runtime.GOOS == "windows" {
		return []string{}
	}
	if _, err := execlog.LookPath("fish"); err != nil {
		return []string{}
	}

	return cache("fish", func() ([]string, error) {
		configDir, err := xdg.UserConfigDir()
		if err != nil {
			return nil, err
		}

		configPath := filepath.Join(configDir, "carapace/bridge/fish/config.fish")
		snippet := fmt.Sprintf(`set __fish_config_dir "/var/empty";source "$__fish_data_dir/config.fish";source %#v;echo $fish_complete_path`, configPath)

		output, err := execlog.Command("fish", "--no-config", "--command", snippet).Output()
		if err != nil {
			return nil, err
		}

		unique := make(map[string]bool)
		for _, location := range strings.Split(string(output), " ") {
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
