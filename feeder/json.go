package feeder

import (
	"encoding/json"
	"io/ioutil"
)

type Json struct {
	Path string
}

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
