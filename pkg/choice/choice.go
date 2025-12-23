package choice

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/carapace-sh/carapace/pkg/xdg"
)

type Choice struct {
	Name    string
	Variant string
	Group   string
}

func (c Choice) Format() string {
	nameVariant := strings.Join([]string{c.Name, c.Variant}, "/")
	return strings.Join([]string{nameVariant, c.Group}, "@")
}

func Parse(s string) Choice {
	nameVariant, group, _ := strings.Cut(s, "@")
	cName, variant, _ := strings.Cut(nameVariant, "/")
	return Choice{
		Name:    cName,
		Variant: variant,
		Group:   group,
	}
}

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
		return errors.New("invalid name")
	}

	if c.Variant == "" || c.Group == "" {
		return errors.New("not fully qualified choice") // TODO force fully qualified for now (instead of git@bridge)
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

func List() ([]string, error) {
	configDir, err := xdg.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, "carapace", "choices")
	entries, err := os.ReadDir(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	choices := make([]string, 0)
	for _, entry := range entries {
		choices = append(choices, entry.Name())
	}
	return choices, nil
}
