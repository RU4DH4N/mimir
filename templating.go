package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// this isn't dangerous at allll........
func renderMarkdown(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read markdown: %w", err)
	}

	var buf bytes.Buffer
	if err := goldmark.Convert(content, &buf); err != nil {
		return "", fmt.Errorf("failed to convert markdown: %w", err)
	}

	return buf.String(), nil
}
