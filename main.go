package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RU4DH4N/mimir/handler"
	"github.com/RU4DH4N/mimir/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func readIndexJson(file string) (Index, error) {
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

func getRoute(title string) string {
	return filepath.Clean("/" + helper.Slugify(title))
}

func getRoutes(root, path string, rMap map[string][]handler.Route) {
	indexPath := filepath.Join(root, path, "index.json")
	index, err := readIndexJson(indexPath)
	if err != nil {
		return
	}

	// Register home pages
	for _, language := range index.Home.Languages {
		route := getRoute(language.Title)
		actual := filepath.Join(root, path, language.File)
		rMap[route] = append(rMap[route], handler.Route{
			Actual: actual,
			Toc:    false,
		})

	}

	for _, page := range index.Pages {
		for _, language := range page.Languages {
			route := getRoute(language.Title)
			actual := filepath.Join(root, path, language.File)
			rMap[route] = append(rMap[route], handler.Route{
				Actual: actual,
				Toc:    true,
			})
		}
	}

	for _, category := range index.SubCategories {
		getRoutes(root, filepath.Join(path, category), rMap)
	}
}

func main() {
	// load 'dir' from config at some point
	dir := "wiki-example"

	e := echo.New()

	e.Use(middleware.Secure())

	e.Static("/static", filepath.Join(dir, "static"))

	// Load Templates
	tmplRoot := filepath.Join(dir, "internal", "templates")
	tmpls := helper.LoadTemplates(tmplRoot)
	e.Renderer = &helper.TemplateRegistry{
		Templates: tmpls,
	}

	w := e.Group("/wiki")

	// load 'en' from config at some point
	contentRoot := filepath.Join(dir, "content")
	rMap := make(map[string][]handler.Route)
	getRoutes(contentRoot, "", rMap)
	for route, actuals := range rMap {
		if len(actuals) == 1 {
			w.GET(route, handler.PageHandler(actuals[0]))
		} else if len(actuals) > 1 {
			title := strings.ReplaceAll(route[1:], "_", " ")

			w.GET(route, handler.DisambiguationHandler(title, contentRoot, actuals))
			//for _, actual := range actuals {
			// not decided how to do this yet
			//}
		}
	}

	// Remove this later
	fmt.Println("=== Registered Routes ===")
	for _, route := range e.Routes() {
		fmt.Printf("[%s] %s -> %s\n", route.Method, route.Path, route.Name)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
