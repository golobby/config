package feeder

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Json is a feeder that feeds using a single json file.
type Json struct {
	Path string
}

// Feed returns all the content.
func (j Json) Feed() (map[string]interface{}, error) {
	fl, err := os.Open(filepath.Clean(j.Path))
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	items := map[string]interface{}{}

	if err := json.NewDecoder(fl).Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}
