package filler

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Json struct {
	Path string
}

func (j Json) Fill(structure interface{}) error {
	file, err := os.Open(filepath.Clean(j.Path))
	if err != nil {
		return err
	}

	if err = json.NewDecoder(file).Decode(structure); err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return err
	}

	return nil
}
