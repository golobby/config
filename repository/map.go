package repository

type Map map[string]interface{}

func (m Map) Extract() (map[string]interface{}, error) {
	return m, nil
}
