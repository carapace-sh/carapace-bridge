package main

import (
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
		switch cmd.Flag("shell").Value.String() {
		case "bash":
			invokeBash(args[0])
		case "elvish":
			invokeElvish(args[0])
		case "fish":
			invokeFish(args[0])
		case "powershell":
			invokePowershell(args[0])
		case "xonsh":
			invokeXonsh(args[0])
		case "zsh":
			invokeZsh(args[0])
		default:

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
	rootCmd.Flags().StringP("format", "f", "json", "shell to use")
	rootCmd.Flags().StringP("shell", "s", "", "shell to use")

	carapace.Gen(rootCmd).FlagCompletion(carapace.ActionMap{
		"format": carapace.ActionValues("json", "tab", "value", "display"),
		"shell":  carapace.ActionValues("bash", "elvish", "fish", "powershell"),
	})
}

type rawValue struct {
	Value       string `json:"value"`
	Display     string `json:"display"`
	Description string `json:"description"`
}

func invokeBash(cmdline string) {
	output, err := exec.Command("scripts/invoke_bash", cmdline).Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(output), "\n")
	vals := make([]*rawValue, 0)
	for _, line := range lines[:len(lines)-1] {
		vals = append(vals, &rawValue{Value: line})
	}
	marshalled, err := json.Marshal(vals)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(marshalled))
}

func invokeElvish(cmdline string) {
	e, _, err := expect.Spawn("elvish", -1)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()

	file, err := ioutil.TempFile(os.TempDir(), "invoke-completion_elvish")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	e.Send(fmt.Sprintf("$edit:completion:arg-completer[%v] %v'' | to-json > %v\n", strings.SplitN(cmdline, " ", 2)[0], cmdline, file.Name()))
	e.Send("echo EXPECT_END\n")
	e.Expect(regexp.MustCompile("EXPECT_END"), 10*time.Second)
	e.Send("exit\n")
	content, err := ioutil.ReadFile(file.Name())
	fmt.Println(string(content))
}

func invokeFish(cmdline string) {
	output, err := exec.Command("scripts/invoke_fish", cmdline).Output()
	if err != nil {
		log.Fatal(err)
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
	marshalled, err := json.Marshal(vals)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(marshalled))
}

func invokeXonsh(cmdline string) {
	e, _, err := expect.SpawnWithArgs([]string{"xonsh", "-i", "--shell-type", "dumb"}, -1)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()

	file, err := ioutil.TempFile(os.TempDir(), "invoke-completion_xonsh")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	e.Send(fmt.Sprintf(`for (k,v) in builtins.__xonsh__.completers.items():
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
	fmt.Println(string(content))
}

type completionResult struct {
	CompletionText string
	ListItemText   string
	ResultType     int
	ToolTip        string
}

func invokePowershell(cmdline string) {
	output, err := exec.Command("scripts/invoke_powershell", cmdline).Output()
	if err != nil {
		log.Fatal(err)
	}

	if !strings.HasPrefix(string(output), "[") {
		output = []byte("[" + string(output) + "]")
	}

	var completionResults []completionResult
	if err := json.Unmarshal(output, &completionResults); err != nil {
		log.Fatal(err.Error())
	}

	vals := make([]*rawValue, 0)
	for _, c := range completionResults {
		vals = append(vals, &rawValue{
			Value:       c.CompletionText,
			Display:     c.ListItemText,
			Description: c.ToolTip,
		})
	}

	marshalled, err := json.Marshal(vals)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(marshalled))
}

func invokeZsh(cmdline string) {
	output, err := exec.Command("scripts/invoke_zsh", cmdline).Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(output), "\n")
	vals := make([]*rawValue, 0)
	r := regexp.MustCompile(`^(?P<value>.*?)( --( (?P<display>.*?))? +-- (?P<description>.*))?$`)
	for _, line := range lines[:len(lines)-1] {
		if r.MatchString(line) {
			matches := r.FindStringSubmatch(line)
			vals = append(vals, &rawValue{Value: matches[1], Display: matches[4], Description: matches[5]})
		}
	}
	marshalled, err := json.Marshal(vals)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(marshalled))
}
