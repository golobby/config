package feeder

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// YamlDirectory is a feeder that feeds using a directory of yaml files.
type YamlDirectory struct {
	Path string
}

// Feed returns all the content.
func (yd YamlDirectory) Feed() (map[string]interface{}, error) {
	files, err := ioutil.ReadDir(yd.Path)
	if err != nil {
		return nil, err
	}

	all := map[string]interface{}{}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		j := Yaml{Path: filepath.Join(yd.Path, string(filepath.Separator), f.Name())}

		items, err := j.Feed()
		if err != nil {
			return nil, err
		}

		k := strings.Split(f.Name(), ".")[0]
		all[k] = items
	}

	return all, nil
}
