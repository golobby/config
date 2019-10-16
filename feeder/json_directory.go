package feeder

import (
	"io/ioutil"
	"os"
	"strings"
)

// JsonDirectory is a feeder that feeds using a directory of json files.
type JsonDirectory struct {
	Path string
}

// Feed will return the feed
func (jd JsonDirectory) Feed() (map[string]interface{}, error) {
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

		items, err := j.Feed()
		if err != nil {
			return nil, err
		}

		k := strings.Split(f.Name(), ".")[0]
		all[k] = items
	}

	return all, nil
}
