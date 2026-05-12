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

		fishConfigDir := filepath.Join(configDir, "carapace/bridge/fish")
		fishConfigFile := filepath.Join(fishConfigDir, "config.fish")
		fishCompletionDir := filepath.Join(configDir, "fish/completions")
		// TODO explicitly adding $__fish_data_dir/completions which is currently missing in $fish_complete_path
		snippet := fmt.Sprintf(`set __fish_config_dir %[1]q;test -f "$__fish_data_dir/config.fish";and source "$__fish_data_dir/config.fish";test -f %[2]q;and source %[2]q;echo %[3]q $fish_complete_path $__fish_data_dir/completions`, fishConfigDir, fishConfigFile, fishCompletionDir)

		output, err := execlog.Command("fish", "--no-config", "--command", snippet).Output()
		if err != nil {
			return nil, err
		}

		unique := make(map[string]bool)
		for _, location := range strings.Fields(string(output)) {
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
