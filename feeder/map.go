package feeder

type Map map[string]interface{}

func (m Map) Feed() (map[string]interface{}, error) {
	return m, nil
}
