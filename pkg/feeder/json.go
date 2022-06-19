package feeder

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

// Json is a feeder.
// It feeds using a JSON file.
type Json struct {
    Path string
}

func (f Json) Feed(structure interface{}) error {
    file, err := os.Open(filepath.Clean(f.Path))
    if err != nil {
        return fmt.Errorf("json: %v", err)
    }

    if err = json.NewDecoder(file).Decode(structure); err != nil {
        return fmt.Errorf("json: %v", err)
    }

    return file.Close()
}
