package bridges

import (
	"github.com/rsteube/carapace-bridge/pkg/env"
)

func filter(m map[string]bool, filter ...[]string) map[string]bool {
	for _, f := range filter {
		for _, e := range f {
			delete(m, e)
		}
	}
	return m
}

func Bridges() map[string]string {
	m := Config()
	for _, bridge := range env.Bridges() {
		if f, ok := map[string]func() []string{
			"bash":         Bash,
			"fish":         Fish,
			"inshellisens": Inshellisense,
			"zsh":          Zsh,
		}[bridge]; ok {
			for _, name := range f() {
				if _, ok := m[name]; !ok {
					m[name] = bridge
				}
			}
		}
	}
	return m
}
