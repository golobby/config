package feeder

// Map is a feeder that feeds using a simple map[string]interface{}
type Map map[string]interface{}

// Feed will return the feed
func (m Map) Feed() (map[string]interface{}, error) {
	return m, nil
}
