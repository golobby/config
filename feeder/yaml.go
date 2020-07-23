package feeder

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Yaml struct {
	Path string
}

func (y *Yaml) Feed() (map[string]interface{}, error) {
	bs, err := ioutil.ReadFile(y.Path)
	if err != nil {
	    return nil, err
	}
	items := make(map[string]interface{})

	err = yaml.Unmarshal(bs, items)
	if err != nil {
	    return nil, err
	}
	return items, nil
}
