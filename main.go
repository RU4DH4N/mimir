package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RU4DH4N/mimir/handler"
	"github.com/RU4DH4N/mimir/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Index struct {
	Home struct {
		Lang []string
		File string
	}
	Pages []struct {
		Lang []string
		Slug string
		File string
	}
	SubCategories []string
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

func getRoutes(root, path, defaultLanguage string, routeMap map[string][]string) {
	indexPath := filepath.Join(root, path, "index.json")
	index, err := readIndexJson(indexPath)
	if err != nil {
		return
	}

	// Register home pages
	for _, lang := range index.Home.Lang {
		var route string
		if lang != defaultLanguage {
			route = "/" + path + "/" + lang
		} else {
			route = "/" + path
		}

		route = filepath.Clean(route) // this isn't great

		actual := filepath.Join(root, path, lang, index.Home.File)
		routeMap[route] = append(routeMap[route], actual)
	}

	// Register individual pages
	for _, page := range index.Pages {
		for _, lang := range page.Lang {
			var route string
			if lang != defaultLanguage {
				route = "/" + lang + page.Slug
			} else {
				route = "/" + page.Slug
			}

			route = filepath.Clean(route)

			actual := filepath.Join(root, path, lang, page.File)
			routeMap[route] = append(routeMap[route], actual)
		}
	}

	// Recurse into subcategories
	for _, sub := range index.SubCategories {
		subPath := filepath.Join(path, sub)
		getRoutes(root, subPath, defaultLanguage, routeMap)
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
	rMap := make(map[string][]string)
	getRoutes(filepath.Join(dir, "content"), "", "en", rMap)
	for route, actuals := range rMap {
		if len(actuals) == 1 {
			w.GET(route, handler.PageHandler(actuals[0]))
		} else if len(actuals) > 1 {
			w.GET(route, handler.DisambiguationHandler(route, filepath.Join(dir, "content"), actuals))
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
