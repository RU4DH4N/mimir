package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/RU4DH4N/mimir/handler"
	"github.com/RU4DH4N/mimir/helper"
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

func (wiki *Wiki) GetRoute(node *Node) (string, string) {
	var actual, route string
	if filepath.Ext(node.Name) == ".md" {
		actual = filepath.Join(wiki.ContentRoot, node.Path)
		route = strings.TrimSuffix(node.Name, ".md")
	} else {
		return "", ""
	}
	return actual, route
}

func (wiki *Wiki) RegisterRoutes(e *echo.Echo) {
	wiki.mu.RLock()
	defer wiki.mu.RUnlock()

	routes := make(map[string][]string)
	wiki.Walk(wiki.ContentTree, func(node *Node) {
		actual, route := wiki.GetRoute(node)
		if actual == "" || route == "" {
			return
		} else if route == "index" { // this wont work
			route = ""
		}
		route = "/" + route
		routes[route] = append(routes[route], actual)
	})

	for route, actuals := range routes {
		amount := len(actuals)
		if amount == 1 {
			e.GET(route, handler.PageHandler(actuals[0]))
		} else if amount > 1 {
			// Handle disambiguation here
			e.GET(route, handler.DisambiguationHandler(route, actuals))
		}
	}

	fmt.Println("=== Registered Routes ===")
	for _, route := range e.Routes() {
		fmt.Printf("[%s] %s -> %s\n", route.Method, route.Path, route.Name)
	}
}

func (wiki *Wiki) Walk(node *Node, fn func(*Node)) {
	if node == nil {
		return
	}
	fn(node)
	for _, child := range node.Children {
		wiki.Walk(child, fn)
	}
}

func main() {
	dir := "wiki-example"
	wiki := LoadWiki(filepath.Join(dir, "content"))

	e := echo.New()
	e.Static("/static", filepath.Join(dir, "static"))

	tmplRoot := filepath.Join(dir, "internal", "templates")
	tmpls := helper.LoadTemplates(tmplRoot)
	e.Renderer = &helper.TemplateRegistry{
		Templates: tmpls,
	}

	wiki.RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
