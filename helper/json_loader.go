package helper

import (
	"encoding/json"
	"fmt"
	"os"
)

func ParseJson(file string, obj any) (any, error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return obj, fmt.Errorf("%s does not exist: %w", file, err)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return obj, fmt.Errorf("error reading %s: %w", file, err)
	}

	if err := json.Unmarshal(data, obj); err != nil {
		return obj, fmt.Errorf("error unmarshaling %s: %w", file, err)
	}

	return obj, nil
}
