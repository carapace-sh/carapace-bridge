package bridges

import (
	"os"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/xdg"
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
	carapace.LOG.Printf("%#v", completers)
	return completers
}
