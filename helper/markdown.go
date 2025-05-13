package helper

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func RenderMarkdown(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to read markdown: %w", err)
	}

	content, err := os.ReadFile(abs)
	if err != nil {
		return "", fmt.Errorf("failed to read markdown: %w", err)
	}

	markdown := goldmark.New( // move this outside the function?
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := markdown.Convert(content, &buf); err != nil {
		return "", fmt.Errorf("failed to convert markdown: %w", err)
	}

	return buf.String(), nil
}
