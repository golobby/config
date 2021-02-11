package feeder

import "os"

// OS is a feeder that feeds using a OS variables.
type OS struct {
	Keys []string
}

// Feed returns all the content.
func (s *OS) Feed() (map[string]interface{}, error) {
	m := make(map[string]interface{})

	if s.Keys != nil {
		for _, k := range s.Keys {
			v := os.Getenv(k)
			if v != "" {
				m[standardize(k)] = v
			}
		}
	}

	return m, nil
}
