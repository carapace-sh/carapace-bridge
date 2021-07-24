package main

import (
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

func invokeBash(cmdline string) {
	output, err := exec.Command("scripts/invoke_bash", cmdline).Output()
	if err != nil {
		log.Fatal(err)
	}
	println(string(output))
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
	println(string(content))
}

func invokeFish(cmdline string) {
	output, err := exec.Command("scripts/invoke_fish", cmdline).Output()
	if err != nil {
		log.Fatal(err)
	}
	println(string(output))
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
    out, _, _ :=e.Expect(regexp.MustCompile("EXPECT_END"), 10*time.Second)
    println(out)
	e.Send("exit\n")
	content, err := ioutil.ReadFile(file.Name())
	println(string(content))
}

func invokePowershell(cmdline string) {
	output, err := exec.Command("scripts/invoke_powershell", cmdline).Output()
	if err != nil {
		log.Fatal(err)
	}
	println(string(output))
}
