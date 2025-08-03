package bridges

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/execlog"
	"github.com/carapace-sh/carapace/pkg/xdg"
)

//go:embed zsh.sh
var zshScript string

func Zsh() []string {
	if runtime.GOOS == "windows" {
		return []string{}
	}

	if _, err := execlog.LookPath("zsh"); err != nil {
		return []string{}
	}

	return cache("zsh", func() ([]string, error) {
		script := zshScript
		if path, err := zshrc(); err == nil {
			script = fmt.Sprintf("autoload -U compinit && compinit;source %#v;%v;compinit", path, script)
		}

		var stdout, stderr bytes.Buffer
		cmd := execlog.Command("zsh", "--no-rcs", "-e", "-c", script)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			if stderr.Len() > 0 {
				carapace.LOG.Println(stderr.String())
			}
			return nil, err
		}

		lines := strings.Split(stdout.String(), "\n")
		unique := make(map[string]bool)
		for _, line := range lines {
			if line != "" {
				unique[line] = true
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

func zshrc() (string, error) {
	configDir, err := xdg.UserConfigDir()
	if err != nil {
		return "", err
	}
	path := configDir + "/carapace/bridge/zsh/.zshrc"
	if _, err := os.Stat(path); err != nil {
		return "", err
	}
	return path, nil
}
