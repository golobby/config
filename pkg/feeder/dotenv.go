package feeder

import (
	"fmt"
	"github.com/golobby/dotenv"
	"os"
	"path/filepath"
)

type DotEnv struct {
	Path string
}

func (f DotEnv) Feed(structure interface{}) error {
	file, err := os.Open(filepath.Clean(f.Path))
	if err != nil {
		return fmt.Errorf("config: cannot open json file; err: %v", err)
	}

	if err = dotenv.NewDecoder(file).Decode(structure); err != nil {
		return fmt.Errorf("config: cannot feed struct; err: %v", err)
	}

	return file.Close()
}
