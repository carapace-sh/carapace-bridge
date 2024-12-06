package bridges

import (
	"os"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/xdg"
	"gopkg.in/yaml.v3"
)

func Config() map[string]string {
	configDir, err := xdg.UserConfigDir()
	if err != nil {
		carapace.LOG.Println(err.Error())
		return make(map[string]string)
	}

	content, err := os.ReadFile(configDir + "/carapace/bridges.yaml")
	if err != nil {
		if !os.IsNotExist(err) {
			carapace.LOG.Println(err.Error())
		}
		return make(map[string]string)
	}

	var completers map[string]string
	if err := yaml.Unmarshal(content, &completers); err != nil {
		carapace.LOG.Println(err.Error())
		return make(map[string]string)
	}

	if completers == nil {
		return make(map[string]string)
	}
	return completers
}
