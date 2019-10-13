package repository

import (
	"io/ioutil"
	"os"
	"strings"
)

type JsonDirectory struct {
	Path string
}

func (jd JsonDirectory) Extract() (map[string]interface{}, error) {
	files, err := ioutil.ReadDir(jd.Path)
	if err != nil {
		return nil, err
	}

	all := map[string]interface{}{}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		j := Json{Path: jd.Path + string(os.PathSeparator) + f.Name()}

		items, err := j.Extract()
		if err != nil {
			return nil, err
		}

		k := strings.Split(f.Name(), ".")[0]
		all[k] = items
	}

	return all, nil
}
