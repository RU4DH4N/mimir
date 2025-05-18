package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/RU4DH4N/mimir/handler"
	"github.com/RU4DH4N/mimir/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func getCategories(root, path string, categories []helper.Category) []helper.Category {
	indexPath := filepath.Join(root, path, "index.json")
	index, err := helper.ParseIndex(indexPath)
	if err != nil {
		return categories
	}

	categories = append(categories, helper.Category{
		Path:  path,
		Home:  &index.Home,
		Pages: index.Pages,
	})

	for _, category := range index.SubCategories {
		categories = getCategories(root, filepath.Join(path, category), categories)
	}
	return categories
}

func loadRoutes(root string, w *echo.Group) error {
	var categories []helper.Category
	categories = getCategories(root, "", categories)
	routes := make(map[string][]handler.PageData)
	for _, category := range categories {
		// load home page for route
		if category.Home != nil && len(category.Home.Languages) > 0 {
			// I'll refactor this to handle languages at some point
			file := filepath.Join(root, category.Path, category.Home.Languages[0].File)
			data := handler.PageData{
				Path:  filepath.Join(category.Path, file),
				Table: false,
			}
			routes[category.Path] = append(routes[category.Path], data)
		}

		for _, page := range category.Pages {
			for _, language := range page.Languages {
				file := filepath.Join(root, category.Path, language.File)
				route := helper.Slugify(language.Title)
				data := handler.PageData{
					Title: language.Title,
					Path:  file,
					Table: true,
				}
				routes[route] = append(routes[route], data)
			}
		}
	}

	// actually register the routes
	for route, data := range routes {
		if len(data) == 1 {
			w.GET(route, handler.PageHandler(data[0]))
		} else if len(data) > 1 {
			w.GET(route, handler.DisambiguationHandler(data[0].Title, root, route, data))
			for i, page := range data {
				r := route + "/v" + strconv.Itoa(i+1)
				w.GET(r, handler.PageHandler(page))
			}
		}
	}

	return nil
}

func main() {
	cfg, err := helper.GetConfig()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		ContentSecurityPolicy: "default-src 'self'",
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// change this
	e.Static("/static", filepath.Join(cfg.WikiRoot, "static"))

	// Load Templates
	tmplRoot := filepath.Join(cfg.WikiRoot, "internal", "templates")
	tmpls := helper.LoadTemplates(tmplRoot)
	e.Renderer = &helper.TemplateRegistry{
		Templates: tmpls,
	}

	e.GET("/wiki", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/wiki/")
	})

	w := e.Group("/wiki/")
	loadRoutes(filepath.Join(cfg.WikiRoot, "content"), w)

	// Remove this later
	fmt.Println("=== Registered Routes ===")
	for _, route := range e.Routes() {
		fmt.Printf("[%s] %s -> %s\n", route.Method, route.Path, route.Name)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	e.Logger.Fatal(e.Start(addr))
}
