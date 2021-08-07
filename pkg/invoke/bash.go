package invoke

import (
	_ "embed"
	"os/exec"
	"strings"
)

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

