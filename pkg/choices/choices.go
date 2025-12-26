package choices

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/carapace-sh/carapace/pkg/xdg"
)

func Get(name string) (*Choice, error) {
	if name == "" || strings.Contains(name, "/") {
		return nil, errors.New("invalid name")
	}

	configDir, err := xdg.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, "carapace", "choices", name)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, errors.New("file is empty")
	}

	choice := Parse(scanner.Text())
	if choice.Name != name {
		return nil, errors.New("invalid content")
	}
	return &choice, nil
}

func Set(c Choice) error {
	if c.Name == "" || strings.Contains(c.Name, "/") {
		return fmt.Errorf("invalid name: %q", c.Format())
	}

	if c.Variant == "" && c.Group == "" {
		return fmt.Errorf("at least one of variant or group must be set: %q", c.Format())
	}

	configDir, err := xdg.UserConfigDir()
	if err != nil {
		return err
	}

	path := filepath.Join(configDir, "carapace", "choices", c.Name)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(c.Format()+"\n"), 0644)
}

func Unset(name string) error {
	if name == "" || strings.Contains(name, "/") {
		return errors.New("invalid name")
	}

	configDir, err := xdg.UserConfigDir()
	if err != nil {
		return err
	}

	path := filepath.Join(configDir, "carapace", "choices", name)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func List(full bool) ([]*Choice, error) {
	configDir, err := xdg.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, "carapace", "choices")
	entries, err := os.ReadDir(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	choices := make([]*Choice, 0)
	for _, entry := range entries {
		switch full {
		case true:
			// TODO slow
			choice, err := Get(entry.Name())
			if err != nil {
				return nil, err
			}
			choices = append(choices, choice)
		default:
			choices = append(choices, &Choice{Name: entry.Name()})
		}
	}
	return choices, nil
}
