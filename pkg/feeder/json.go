package feeder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Json struct {
	Path string
}

func (f Json) Feed(structure interface{}) error {
	file, err := os.Open(filepath.Clean(f.Path))
	if err != nil {
		return fmt.Errorf("config: cannot open json file; err: %v", err)
	}

	if err = json.NewDecoder(file).Decode(structure); err != nil {
		return fmt.Errorf("config: cannot feed struct; err: %v", err)
	}

	return file.Close()
}
