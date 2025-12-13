package bridge

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/carapace-sh/carapace"
	shlex "github.com/carapace-sh/carapace-shlex"
	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/carapace-sh/carapace/pkg/xdg"
)

func ensureExists(path string) (err error) {
	if _, err = os.Stat(path); err == nil || !os.IsNotExist(err) {
		return
	}
	if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return
	}
	_, err = os.Create(path)
	return
}

// ActionPowershell bridges completions registered in powershell
//
//	(uses custom `Microsoft.PowerShell_profile.ps1` in `~/.config/carapace/bridge/powershell`)
func ActionPowershell(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			configDir, err := xdg.UserConfigDir()
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}
			configPath := fmt.Sprintf("%v/carapace/bridge/powershell/Microsoft.PowerShell_profile.ps1", configDir)
			if err := ensureExists(configPath); err != nil {
				return carapace.ActionMessage(err.Error())
			}

			args := append(command, c.Args...)
			args = append(args, c.Value)

			// for index, arg := range args {
			// TODO handle different escape character and escaping in general
			// args[index] = strings.Replace(arg, " ", "` ", -1)
			// }

			line := shlex.Join(args)
			snippet := []string{
				fmt.Sprintf(`Get-Content "%v/carapace/bridge/powershell/Microsoft.PowerShell_profile.ps1" | Out-String | Invoke-Expression`, configDir),
				fmt.Sprintf(`[System.Management.Automation.CommandCompletion]::CompleteInput("%v", %v, $null).CompletionMatches | ConvertTo-Json `, line, len(line)),
			}
			return carapace.ActionExecCommand("pwsh", "-Command", strings.Join(snippet, ";"))(func(output []byte) carapace.Action {
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
				a := carapace.ActionValuesDescribed(vals...).StyleF(style.ForPath)
				if len(suffixes) > 0 {
					return a.NoSpace(suffixes...)
				}
				return a
			}).Invoke(c).ToA()
		})
	})
}
