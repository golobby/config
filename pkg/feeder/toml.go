package feeder

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Yaml is a feeder.
// It feeds using a YAML file.
type Toml struct {
	Path string
}

func (f Toml) Feed(structure interface{}) error {
	tomlContent, err := ioutil.ReadFile(filepath.Clean(f.Path))
	if err != nil {
		return fmt.Errorf("config: cannot read toml file; err: %v", err)
	}

	tomlString := string(tomlContent)

	if _, err = toml.Decode(tomlString, structure); err != nil {
		return fmt.Errorf("config: cannot feed struct; err: %v", err)
	}

	return nil
}
