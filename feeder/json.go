package feeder

import (
	"encoding/json"
	"io/ioutil"
)

// Json is a feeder that feeds using a single json file.
type Json struct {
	Path string
}

// Feed returns all the content.
func (j Json) Feed() (map[string]interface{}, error) {
	content, err := ioutil.ReadFile(j.Path)
	if err != nil {
		return nil, err
	}

	items := map[string]interface{}{}

	if err := json.Unmarshal(content, &items); err != nil {
		return nil, err
	}

	return items, nil
}
