// Package feeder is a collection of feeders
package feeder

// Map is a feeder that feeds using a simple runtime map[string]interface{}
type Map map[string]interface{}

// Feed returns all the content.
func (m Map) Feed() (map[string]interface{}, error) {
	return m, nil
}
