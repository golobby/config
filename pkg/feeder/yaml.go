package feeder

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Yaml struct {
	Path string
}

func (f Yaml) Feed(structure interface{}) error {
	file, err := os.Open(filepath.Clean(f.Path))
	if err != nil {
		return fmt.Errorf("config: cannot open json file; err: %v", err)
	}

	if err = yaml.NewDecoder(file).Decode(structure); err != nil {
		return fmt.Errorf("config: cannot feed struct; err: %v", err)
	}

	return file.Close()
}
