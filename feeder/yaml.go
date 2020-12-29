package feeder

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Yaml struct {
	Path string
}

func (y *Yaml) Feed() (map[string]interface{}, error) {
	fl, err := os.Open(filepath.Clean(y.Path))
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	items := make(map[string]interface{})

	if err := yaml.NewDecoder(fl).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}
