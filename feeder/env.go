package feeder

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Env is a feeder that feeds using a single environment file.
type Env struct {
	Path               string
	DisableOSVariables bool
}

// Feed returns all the content.
func (e *Env) Feed() (map[string]interface{}, error) {
	m := make(map[string]interface{})

	values, err := e.load(e.Path)
	if err != nil {
		return nil, err
	}

	for k, v := range values {
		m[standardize(k)] = e.get(k, v)
	}

	return m, nil
}

// load reads the given env file and extracts the variables as a string map
func (e *Env) load(filename string) (map[string]string, error) {
	path, _ := filepath.Abs(filename)
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	variables, err := e.read(file)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

// read opens the given env file and extracts the variables as a string map
func (e *Env) read(file io.Reader) (map[string]string, error) {
	items := map[string]string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		key, value, err := e.parse(scanner.Text())
		if err != nil {
			return nil, err
		}

		if key != "" {
			items[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// parse extracts the key/value from the given line
func (e *Env) parse(line string) (string, string, error) {
	ln := strings.TrimSpace(line)

	if len(ln) == 0 {
		return "", "", nil
	}

	if ln[0] == '#' {
		return "", "", nil
	}

	s := strings.Index(ln, "=")
	if s == -1 {
		return "", "", errors.New("Invalid line: " + ln)
	}

	k := strings.TrimSpace(ln[:s])
	v := strings.TrimSpace(ln[s+1:])

	return k, v, nil
}

// get fetches variable from OS if exist and return fallback otherwise.
func (e *Env) get(key, fallback string) string {
	if !e.DisableOSVariables {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}

	return fallback
}

// standardize updates config key (e.g. APP_NAME to  app.name)
func standardize(k string) string {
	return strings.Replace(strings.ToLower(k), "_", ".", -1)
}
