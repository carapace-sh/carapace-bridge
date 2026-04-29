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
			for _, name := range expandZshCompdef(line) {
				if name != "" {
					unique[name] = true
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

func expandZshCompdef(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" || strings.HasPrefix(s, "_") || strings.HasPrefix(s, "-") || strings.Contains(s, "=") {
		return nil
	}

	if strings.HasPrefix(s, "(") {
		if end := strings.Index(s, ")"); end > 0 {
			var expanded []string
			for _, alternative := range strings.Split(s[1:end], "|") {
				for _, prefix := range expandCharClasses(alternative) {
					expanded = append(expanded, trimZshPattern(prefix+s[end+1:]))
				}
			}
			return filterNames(expanded)
		}
	}

	return filterNames([]string{trimZshPattern(s)})
}

func expandCharClasses(s string) []string {
	start := strings.Index(s, "[")
	if start == -1 {
		return []string{s}
	}
	end := strings.Index(s[start:], "]")
	if end == -1 {
		return []string{s}
	}
	end += start

	var expanded []string
	for _, r := range s[start+1 : end] {
		for _, suffix := range expandCharClasses(s[end+1:]) {
			expanded = append(expanded, s[:start]+string(r)+suffix)
		}
	}
	return expanded
}

func trimZshPattern(s string) string {
	for i, r := range s {
		switch r {
		case '[', ']', '(', ')', '|', '*', '?', '#':
			return s[:i]
		}
	}
	return s
}

func filterNames(names []string) []string {
	var filtered []string
	for _, name := range names {
		if name == "" || len(name) == 1 || strings.ContainsAny(name, "[]*?#()|") {
			continue
		}
		filtered = append(filtered, name)
	}
	return filtered
}
