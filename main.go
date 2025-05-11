package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
)

type Node struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*Node
}

type Wiki struct {
	ContentRoot string
	ContentTree *Node
	mu          sync.RWMutex
}

func LoadWiki(contentRoot string) *Wiki {
	wiki := &Wiki{
		ContentRoot: contentRoot,
	}

	wiki.mu.Lock()
	defer wiki.mu.Unlock()

	if _, err := os.Stat(wiki.ContentRoot); os.IsNotExist(err) {
		log.Fatalf("content root does not exist: %s", wiki.ContentRoot)
	}

	wiki.ContentTree = &Node{
		Name:     filepath.Base(wiki.ContentRoot),
		Path:     "",
		IsDir:    true,
		Children: []*Node{},
	}

	wiki.traverseContent(wiki.ContentRoot, wiki.ContentTree)
	return wiki
}

func (wiki *Wiki) traverseContent(currentPath string, parentNode *Node) {
	items, err := os.ReadDir(currentPath)
	if err != nil {
		log.Printf("Error reading directory %q: %v", currentPath, err)
		return
	}

	for _, item := range items {
		name := item.Name()
		itemPath := filepath.Clean(filepath.Join(currentPath, name))

		relativePath, err := filepath.Rel(wiki.ContentRoot, itemPath)
		if err != nil {
			log.Printf("Error getting relative path for %q: %v", itemPath, err)
			continue
		}

		node := &Node{
			Name:     name,
			Path:     filepath.ToSlash(relativePath),
			IsDir:    item.IsDir(),
			Children: []*Node{},
		}

		parentNode.Children = append(parentNode.Children, node)

		if node.IsDir {
			wiki.traverseContent(itemPath, node)
		}
	}
}

func (wiki *Wiki) RegisterRoutes(e *echo.Echo) {
	wiki.mu.RLock()
	defer wiki.mu.RUnlock()

	nameToPaths := make(map[string][]string)

	var register func(node *Node)
	register = func(node *Node) {
		if node.IsDir {
			indexPath := filepath.Join(wiki.ContentRoot, node.Path, "index.md")
			routePath := "/" + node.Path
			if node.Path == "" {
				routePath = "/"
			}

			if _, err := os.Stat(indexPath); err == nil {
				e.GET(routePath, wiki.handlerFor(indexPath))
			} else {
				e.GET(routePath, func(c echo.Context) error {
					return echo.NewHTTPError(http.StatusNotFound, "Page not found")
				})
			}
		} else if filepath.Ext(node.Name) == ".md" {
			if strings.HasSuffix(node.Name, "index.md") {
				return
			}

			cleanPath := "/" + strings.TrimSuffix(node.Path, ".md")
			fullPath := filepath.Join(wiki.ContentRoot, node.Path)
			e.GET(cleanPath, wiki.handlerFor(fullPath))

			base := strings.TrimSuffix(node.Name, ".md")
			nameToPaths[base] = append(nameToPaths[base], node.Path)
		}

		for _, child := range node.Children {
			register(child)
		}
	}
	register(wiki.ContentTree)

	// Register short names and disambiguation
	for name, paths := range nameToPaths {
		route := "/" + name

		if len(paths) == 1 {
			fullPath := filepath.Join(wiki.ContentRoot, paths[0])
			e.GET(route, wiki.handlerFor(fullPath))
		} else {
			// Copy to avoid closure capture issue
			nameCopy := name
			pathsCopy := append([]string(nil), paths...)

			e.GET(route, wiki.disambiguationHandler(nameCopy, pathsCopy))
		}
	}

	fmt.Println("=== Registered Routes ===")
	for _, route := range e.Routes() {
		fmt.Printf("[%s] %s -> %s\n", route.Method, route.Path, route.Name)
	}
}

func (wiki *Wiki) disambiguationHandler(title string, paths []string) echo.HandlerFunc {
	return func(c echo.Context) error {
		links := make([]string, len(paths))
		for i, p := range paths {
			links[i] = "/" + strings.TrimSuffix(p, ".md")
		}
		data := map[string]interface{}{
			"Title": title,
			"Links": links,
		}
		return c.Render(http.StatusOK, "base.html", data)
	}
}

func (wiki *Wiki) handlerFor(fullPath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		htmlContent, err := renderMarkdown(fullPath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render page")
		}
		data := map[string]interface{}{
			"Content": template.HTML(htmlContent),
		}
		return c.Render(http.StatusOK, "base.html", data)
	}
}

func loadTemplates(root string) *template.Template {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(fmt.Errorf("failed to walk templates: %w", err))
	}

	return template.Must(template.ParseFiles(files...))
}

func main() {
	dir := "wiki-example"
	wiki := LoadWiki(filepath.Join(dir, "content"))

	e := echo.New()
	e.Static("/static", filepath.Join(dir, "static"))

	tmplRoot := filepath.Join(dir, "internal", "templates")
	tmpls := loadTemplates(tmplRoot)
	e.Renderer = &TemplateRenderer{templates: tmpls}

	wiki.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
