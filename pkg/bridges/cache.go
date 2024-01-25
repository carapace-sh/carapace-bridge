package bridges

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/xdg"
)

func cache(s string, f func() ([]string, error)) []string {
	cacheDir, err := xdg.UserCacheDir()
	if err != nil {
		return []string{}
	}

	path := fmt.Sprintf("%v/carapace/bridges-%s.json", cacheDir, s)
	names, err := load(path)
	if err == nil {
		return names
	}
	carapace.LOG.Println(err.Error())

	names, err = f()
	if err != nil {
		carapace.LOG.Println(err.Error())
		return []string{}
	}

	if err := save(path, names); err != nil {
		carapace.LOG.Println(err.Error())
	}

	return names
}

func load(path string) ([]string, error) {
	timeout := 24 * time.Hour
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.ModTime().Add(timeout).Before(time.Now()) {
		return nil, errors.New("timeout exceeded")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var names []string
	if err := json.Unmarshal(content, &names); err != nil {
		return nil, err
	}
	return names, nil
}

func save(path string, names []string) error {
	m, err := json.Marshal(names)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	return os.WriteFile(path, m, os.ModePerm)
}
