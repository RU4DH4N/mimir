package helper

import (
	"fmt"
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
	val, err := ParseJson(file, &idx)
	if err != nil {
		return idx, fmt.Errorf("unable to get index: %w", err)
	}

	idxPtr, ok := val.(*Index)
	if !ok {
		return idx, fmt.Errorf("unable to parse index")
	}

	return *idxPtr, nil
}
