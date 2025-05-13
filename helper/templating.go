package helper

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type TemplateRegistry struct {
	Templates map[string]*template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data any, c echo.Context) error {
	tmpl, ok := t.Templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}

// I know this is crap, I'll rewrite it later
func LoadTemplates(root string) map[string]*template.Template {
	templates := make(map[string]*template.Template)

	var partials []string
	var base string

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}

		name := filepath.Base(path)
		switch name {
		case "base.html":
			base = path
		default:
			partials = append(partials, path)
		}
		return nil
	})

	if base == "" {
		fmt.Println("base.html not found")
		return templates
	}

	for _, path := range partials {
		name := filepath.Base(path)

		files := []string{base}
		for _, p := range partials {
			if p != path {
				files = append(files, p)
			}
		}
		files = append(files, path)

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Printf("failed to parse template %s: %v\n", path, err)
			continue
		}

		templates[name] = tmpl
	}

	return templates
}
