package bridges

import (
	"encoding/json"

	"github.com/carapace-sh/carapace/pkg/execlog"
)

func Inshellisense() []string {
	if _, err := execlog.LookPath("inshellisense"); err != nil {
		return []string{}
	}

	return cache("inshellisense", func() ([]string, error) {
		out, err := execlog.Command("inshellisense", "specs", "list").Output()
		if err != nil {
			return nil, err
		}

		var entries []string
		if err := json.Unmarshal(out, &entries); err != nil {
			return nil, err
		}

		unique := make(map[string]bool)
		for _, entry := range entries {
			unique[entry] = true
		}

		completers := make([]string, 0)
		for c := range filter(unique, inshellisenseBuiltins) {
			completers = append(completers, c)
		}
		return completers, nil
	})
}
