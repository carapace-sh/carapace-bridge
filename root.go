package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	expect "github.com/google/goexpect"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "invoke-completion",
	Short: "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var f func(string) ([]*rawValue, error)
		switch cmd.Flag("shell").Value.String() {
		case "bash":
			f = InvokeBash
		case "elvish":
			f = InvokeElvish
		case "fish":
			f = InvokeFish
		case "oil":
			f = InvokeOil
		case "powershell":
			f = InvokePowershell
		case "xonsh":
			f = InvokeXonsh
		case "zsh":
			f = InvokeZsh
		default:
			log.Fatal("TODO: determine shell") // TODO determine shell
		}

		if vals, err := f(args[0]); err != nil {
			log.Fatal(err.Error())
		} else {
			switch cmd.Flag("format").Value.String() {
			case "json":
				m, err := json.Marshal(vals)
				if err != nil {
					log.Fatal(err.Error())
				}
				fmt.Println(string(m))
			case "value":
				for _, v := range vals {
					fmt.Println(v.Value)
				}
			default:
				log.Fatal("unknown format")
			}
		}
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func main() {
	rootCmd.Execute()
}
func init() {
	rootCmd.Flags().StringP("format", "f", "json", "output format")
	rootCmd.Flags().StringP("shell", "s", "", "shell to use")

	carapace.Gen(rootCmd).FlagCompletion(carapace.ActionMap{
		"format": carapace.ActionValues("json", "value"),
		"shell":  carapace.ActionValues("bash", "elvish", "fish", "oil", "powershell", "xonsh", "zsh"),
	})
}

type rawValue struct {
	Value       string `json:"value"`
	Display     string `json:"display"`
	Description string `json:"description"`
}

//go:embed scripts/invoke_bash
var bashScript string

func InvokeBash(cmdline string) ([]*rawValue, error) {
	output, err := exec.Command("bash", "-i", "-c", bashScript, "--", cmdline).Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	vals := make([]*rawValue, 0)
	for _, line := range lines[:len(lines)-1] {
		vals = append(vals, &rawValue{Value: line})
	}
	return vals, nil
}

type complexCandidate struct {
	Stem       string
	CodeSuffix string
	Display    string
}

func InvokeElvish(cmdline string) ([]*rawValue, error) {
	e, _, err := expect.Spawn("elvish", -1)
	if err != nil {
		return nil, err
	}
	defer e.Close()

	file, err := ioutil.TempFile(os.TempDir(), "invoke-completion_elvish")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())

	e.Send(fmt.Sprintf("$edit:completion:arg-completer[%v] %v'' | put [(all)] | to-json > %v\n", strings.SplitN(cmdline, " ", 2)[0], cmdline, file.Name()))
	e.Send("echo EXPECT_END\n")
	e.Expect(regexp.MustCompile("EXPECT_END"), 10*time.Second)
	e.Send("exit\n")
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return nil, err
	}

	var candidates []complexCandidate
	if err := json.Unmarshal(content, &candidates); err != nil {
		return nil, err
	}

	vals := make([]*rawValue, 0)
	for _, candidate := range candidates {
		vals = append(vals, &rawValue{Value: candidate.Stem + candidate.CodeSuffix, Display: candidate.Display})
	}
	return vals, nil
}

func InvokeFish(cmdline string) ([]*rawValue, error) {
	output, err := exec.Command("fish", "-c", `complete --do-complete="$argv"`, "--", cmdline).Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	vals := make([]*rawValue, 0)
	for _, line := range lines[:len(lines)-1] {
		splitted := strings.SplitN(line, "\t", 2)
		if len(splitted) > 1 {
			vals = append(vals, &rawValue{Value: splitted[0], Description: splitted[1]})
		} else {
			vals = append(vals, &rawValue{Value: splitted[0]})
		}
	}
	return vals, nil
}

// TODO return value
func InvokeXonsh(cmdline string) ([]*rawValue, error) {
	e, _, err := expect.SpawnWithArgs([]string{"xonsh", "-i", "--shell-type", "dumb"}, -1)
	if err != nil {
		return nil, err
	}
	defer e.Close()

	file, err := ioutil.TempFile(os.TempDir(), "invoke-completion_xonsh")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())

	e.Send(fmt.Sprintf(`import builtins
for (k,v) in builtins.__xonsh__.completers.items():
	   e = v('', '%v', 0, len('%v'), '')
	   if e is not None and len(e)!=0:
	       with open('%v', 'a') as f:
                 import json
                 m = list(map(lambda x: dict(value = str(x), display = x.display, description = x.description), e))
                 print(json.dumps(m), file=f)
                 break

`, cmdline, cmdline, file.Name()))
	e.Send("echo EXPECT_END\n")
	e.Expect(regexp.MustCompile("EXPECT_END"), 10*time.Second)
	e.Send("exit\n")
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return nil, err
	}

	var vals []rawValue
	if err := json.Unmarshal(content, &vals); err != nil {
		return nil, err
	}

	// TODO yuck
	result := make([]*rawValue, 0)
	for _, x := range vals {
		result = append(result, &x)
	}
	return result, nil
}

//go:embed scripts/invoke_oil
var oilScript string

func InvokeOil(cmdline string) ([]*rawValue, error) {
	output, err := exec.Command("osh", "-c", oilScript, "--", cmdline).Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	vals := make([]*rawValue, 0)
	for _, line := range lines[:len(lines)-1] {
		vals = append(vals, &rawValue{Value: line})
	}
	return vals, nil
}

type completionResult struct {
	CompletionText string
	ListItemText   string
	ResultType     int
	ToolTip        string
}

func InvokePowershell(cmdline string) ([]*rawValue, error) {
	output, err := exec.Command("pwsh", "-Command", fmt.Sprintf(`[System.Management.Automation.CommandCompletion]::CompleteInput("%v", %v, $null).CompletionMatches | ConvertTo-Json`, cmdline, len(cmdline))).Output()
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(string(output), "[") {
		output = []byte("[" + string(output) + "]")
	}

	var completionResults []completionResult
	if err := json.Unmarshal(output, &completionResults); err != nil {
		return nil, err
	}

	vals := make([]*rawValue, 0)
	for _, c := range completionResults {
		vals = append(vals, &rawValue{
			Value:       c.CompletionText,
			Display:     c.ListItemText,
			Description: c.ToolTip,
		})
	}
	return vals, nil
}

//go:embed scripts/invoke_zsh
var zshScript string

func InvokeZsh(cmdline string) ([]*rawValue, error) {
	output, err := exec.Command("zsh", "-c", zshScript, "--", cmdline).Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\r\n")
	vals := make([]*rawValue, 0)
	rValueOnly := regexp.MustCompile("^(?P<value>.*) -- $")
	r := regexp.MustCompile("^(?P<value>.*?)( --( (?P<display>.*?))?( +-- (?P<description>.*)))?$")
	for _, line := range lines[:len(lines)-1] {
		if rValueOnly.MatchString(line) {
			matches := rValueOnly.FindStringSubmatch(line)
			vals = append(vals, &rawValue{Value: matches[1]})
		} else if r.MatchString(line) {
			matches := r.FindStringSubmatch(line)
			vals = append(vals, &rawValue{Value: matches[1], Display: matches[4], Description: matches[6]})
		}
	}
	return vals, nil
}
