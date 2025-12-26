package bridge

import (
	"os"
	"os/exec"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/x"
	"github.com/spf13/cobra"
)

type bridge struct {
	Name   string
	Action func(command ...string) carapace.Action
}

var candidates = []string{
	// go
	"carapace",
	"cobra",
	"complete",
	"kingpin",
	"urfavecli",
	"urfavecli_v1",

	// python
	"argcomplete",
	"argcomplete_v1",
	"click",

	// rust
	// "clap", // TODO clap dynamic completion is still in development

	// javascript
	"yargs",
}

// TODO experimental
func Detect(cmd string) (*bridge, bool) {
	if _, err := exec.LookPath(cmd); err != nil {
		return nil, false
	}

	for _, candidate := range candidates {
		ok, err := check(bridgeActions[candidate](cmd))
		if err != nil {
			// TODO abort on error?
			carapace.LOG.Println(err.Error())
			continue
		}
		if ok {
			return &bridge{candidate, bridgeActions[candidate]}, true
		}
	}
	return nil, false
}

func check(action carapace.Action) (bool, error) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "carapace-bridge_")
	if err != nil {
		return false, err
	}
	defer os.RemoveAll(tmpDir)

	cmd := &cobra.Command{
		DisableFlagParsing: true,
		Run:                func(cmd *cobra.Command, args []string) {},
	}
	carapace.Gen(cmd).PositionalAnyCompletion(action.Chdir(tmpDir))

	e, err := x.Complete(cmd, "", "", "-")
	if err != nil {
		return false, err
	}

	for _, value := range e.Values {
		switch value.Value {
		case "-h", "-help", "--help":
			return true, nil
		}
	}
	return false, nil
}
