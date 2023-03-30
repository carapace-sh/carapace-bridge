package bridge

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/rsteube/carapace/pkg/xdg"
)

// ActionPowershell bridges completions registered in powershell
// (uses custom `Microsoft.PowerShell_profile.ps1` in “~/.config/carapace/bridge/powershell`)
func ActionPowershell(command ...string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if len(command) == 0 {
			return carapace.ActionMessage("missing argument [ActionPowershell]")
		}

		configDir, err := xdg.UserConfigDir()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}
		c.Setenv("XDG_CONFIG_HOME", fmt.Sprintf("%v/carapace/bridge", configDir))

		args := append(command, c.Args...)
		args = append(args, c.CallbackValue)

		// for index, arg := range args {
		// TODO handle different escape character and escapcing in general
		// args[index] = strings.Replace(arg, " ", "` ", -1)
		// }

		line := strings.Join(args, " ")
		snippet := fmt.Sprintf(`[System.Management.Automation.CommandCompletion]::CompleteInput("%v", %v, $null).CompletionMatches | ConvertTo-Json `, line, len(line))
		return carapace.ActionExecCommand("pwsh", "-Command", snippet)(func(output []byte) carapace.Action {
			if len(output) == 0 {
				return carapace.ActionValues()
			}

			type singleResult struct {
				CompletionText string `json:"CompletionText"`
				ListItemText   string `json:"ListItemText"`
				ResultType     int    `json:"ResultType"`
				ToolTip        string `json:"ToolTip"`
			}
			var result []singleResult

			if err := json.Unmarshal(output, &result); err != nil {
				result = make([]singleResult, 1)
				if err := json.Unmarshal(output, &result[0]); err != nil {
					carapace.LOG.Println(string(output))
					return carapace.ActionMessage(err.Error())
				}
			}

			suffixes := make([]rune, 0)
			vals := make([]string, 0)
			for _, r := range result {
				if _runes := []rune(r.CompletionText); len(_runes) > 2 && strings.HasSuffix(r.CompletionText, " ") {
					suffixes = append(suffixes, _runes[len(_runes)-1])
				}
				r.CompletionText = strings.TrimSuffix(r.CompletionText, " ")

				if r.CompletionText == r.ToolTip {
					r.ToolTip = ""
				}
				vals = append(vals, r.CompletionText, r.ToolTip)
			}
			return carapace.ActionValuesDescribed(vals...).NoSpace(suffixes...).StyleF(style.ForPath)
		}).Invoke(c).ToA()
	})
}
