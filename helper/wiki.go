package helper

import (
	"encoding/json"
	"fmt"
	"os"
)

type Index struct {
	Pages         []Page   `json:"pages"`
	Home          Home     `json:"home"`
	SubCategories []string `json:"subcategories"`
}

type Page struct {
	Languages []Language `json:"languages"`
}

type Home struct {
	Languages []Language `json:"languages"`
}

type Language struct {
	Language string `json:"language"`
	Title    string `json:"title,omitempty"`
	File     string `json:"file"`
}

type Category struct {
	Path  string
	Home  *Home
	Pages []Page
}

func ParseIndex(file string) (Index, error) {
	var idx Index

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return idx, fmt.Errorf("index.json does not exist: %w", err)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return idx, fmt.Errorf("error reading index.json: %w", err)
	}

	if err := json.Unmarshal(data, &idx); err != nil {
		return idx, fmt.Errorf("error unmarshaling index.json: %w", err)
	}

	return idx, nil
}
